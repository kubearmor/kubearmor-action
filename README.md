# KubeArmor Action
A Github Action component library that visualizes the application's system-level behaviors and network connection changes. For example, which processes are generated by the application and the parent-child relationship between processes, which file access is generated, which network connections are generated, network connection topology, network connection changes, and so on.

## How To Use

### Main Action: kubearmor-action
This action will be used to save the new app summary report and choose to generate visualisation results or not
```yaml
 # Save the new app summary report and Choose to Generate visualisation results or not
- name: Save the new app summary report and Choose to Generate visualisation results
  uses: kubearmor/kubearmor-action@main
  with:
    old-summary-path: 'https://raw.githubusercontent.com/kubearmor/kubearmor-action/main/test/testdata/old-summary-data.json'
    namespace: 'sock-shop'
    app-name: 'front-end'
    file: 'summary-test.json'
    install-kubearmor: 'true' # default value is false, if set true, will install kubearmor components
    save-summary-report: 'true' # default value is false, if set true, will save summary report to artifacts
    visualise: 'true' # default value is false, if set true, will generate visualisation results
```

### Complete Example
```yaml
name: test

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
permissions:
    pull-requests: write
    contents: write

jobs:
  test_job:
    runs-on: ubuntu-latest
    name: A job to test kubearmor-action
    steps:
      # Setup Go
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'
        env:
          GOPATH: ${{ runner.workspace }}
          GO111MODULE: "on"
      # Checkout to your repo
      - name: Checkout
        uses: actions/checkout@v3
      # Install k3s cluster(You can setup a k8s cluster here)
      - name: Setup k3s cluster
        run: |
          curl -sfL https://get.k3s.io | K3S_KUBECONFIG_MODE="644" sh -s - --disable=traefik

          KUBEDIR=$HOME/.kube
          KUBECONFIG=$KUBEDIR/config
          [[ ! -d $KUBEDIR ]] && mkdir $KUBEDIR
          if [ -f $KUBECONFIG ]; then
            echo "Found $KUBECONFIG already in place ... backing it up to $KUBECONFIG.backup"
            cp $KUBECONFIG $KUBECONFIG.backup
          fi
          sudo cp /etc/rancher/k3s/k3s.yaml $KUBECONFIG
          sudo chown $USER:$USER $KUBECONFIG
          echo "export KUBECONFIG=$KUBECONFIG" | tee -a ~/.bashrc
          
          echo "wait for initialization"
          sleep 15

          runtime="15 minute"
          endtime=$(date -ud "$runtime" +%s)

          while [[ $(date -u +%s) -le $endtime ]]
          do
            status=$(kubectl get pods -A -o jsonpath={.items[*].status.phase})
            [[ $(echo $status | grep -v Running | wc -l) -eq 0 ]] && break
            echo "wait for initialization"
            sleep 1
          done
          kubectl get pods -A
      # Install kubearmor components(This will install kubearmor-client and Discovery-Engine)
      - name: Install kubearmor components
        uses: kubearmor/kubearmor-action/@main
        with:
          install-kubearmor: 'true'
      # Show pods info
      - name: Get pod
        run: kubectl get po -A
      # Deploy the new app(You can deploy your application here)
      - name: Deploy the new app
        run: kubectl apply -f ./test/testdata/sock-shop.yaml
      # Check all pods are ready, if not, get reason
      - name: Check all pods are ready, if not, get reason
        uses: kubearmor/kubearmor-action/actions/check-pods-ready@main
      # Runs Integration/Tests/Load Generation(You can add a step here)
      # Generate load on the new app
      - name: Generate load on the new app
        run: |
          sleep 60
          docker run --net=host weaveworksdemos/load-test -h localhost:30001 -r 100 -c 2
      # Save the new app summary report and Generate visualisation results
      - name: Save the new app summary report and Generate visualisation results
        uses: kubearmor/kubearmor-action@main
        id: visualisation
        with:
          old-summary-path: 'https://raw.githubusercontent.com/kubearmor/kubearmor-action/gh-pages/latest-summary-test.json'
          namespace: 'sock-shop'
          file: 'latest-summary-test.json'
          save-summary-report: 'true'
          visualise: 'true'
      # Get the latest summary report file
      - name: Get the latest summary report file
        uses: actions/download-artifact@v2
        with:
          name: ${{ steps.visualisation.outputs.summary-report-artifact }}
          path: summary_reports
      # Get the visualisation results
      - name: Get the visualisation results 
        uses: actions/download-artifact@v2
        with:
          name: ${{ steps.visualisation.outputs.visualisation-results-artifact }}
          path: images
      # Store the latest summary report file
      - name: Store the latest summary report file
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./summary_reports
          keep_files: true
      # Store the visualisation results
      - name: Store the visualisation results
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./images
          keep_files: true
      # Comment the visualisation results on the PR
      - name: Comment on PR
        run: |
          gh pr comment  ${{ github.event.number }} -b"![system_graph](https://raw.githubusercontent.com/${{ github.repository }}/gh-pages/${{ steps.visualisation.outputs.sys-visualisation-image }})"
          gh pr comment  ${{ github.event.number }} -b"![network_graph](https://raw.githubusercontent.com/${{ github.repository }}/gh-pages/${{ steps.visualisation.outputs.network-visualisation-image }})"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      # Delete the new app
      - name: Delete the new app
        run: kubectl delete -f ./test/testdata/sock-shop.yaml
```

### Other Tool Actions
#### Action: check-pods-ready
This action will be used to check whether all pods are ready, if not, will show logs and events for troubleshooting.(This will need setup GO env first.)
```yaml
# Check all pods are ready, if not, get reason
- name: Check all pods are ready, if not, get reason
  uses: kubearmor/kubearmor-action/actions/check-pods-ready@main
```

## Application Behaviors Visualisation
### System Behaviors
![Alt text](docs/pics/sys-behaviors.png)
### Network Connections
![Alt text](docs/pics/network-connnections.png)


## Architecture Overview
```Shell
.
├── LICENSE
├── Makefile
├── README.md
├── action.yml
├── actions
│   ├── check-pods-ready
│   │   ├── action.yml
│   │   └── main.go
│   ├── install-kubearmor
│   │   └── action.yml
│   ├── save-summary-report
│   │   └── action.yml
│   ├── setup-k3s-cluster
│   │   └── action.yml
│   └── visual-report
│       └── action.yml
├── cmd
│   └── visual
│       ├── cmd
│       │   ├── common.go
│       │   ├── network.go
│       │   ├── root.go
│       │   └── system.go
│       └── main.go
├── common
│   └── common.go
├── docs
│   └── pics
│       ├── network-connnections.png
│       └── sys-behaviors.png
├── examples
│   └── visualisation
│       └── main.go
├── go.mod
├── go.sum
├── hack
│   ├── LICENSE_TEMPLATE
│   └── install-k3s.sh
├── install
│   ├── k3s
│   │   └── install_k3s.sh
│   └── self-managed-k8s
│       └── crio
│           └── install_crio.sh
├── pkg
│   ├── controller
│   │   └── client
│   │       └── client.go
│   └── visualisation
│       ├── plantuml.jar
│       ├── types.go
│       └── visualisation.go
├── test
│   └── testdata
│       ├── new-summary-data.json
│       ├── old-summary-data.json
│       ├── sock-shop.yaml
│       └── wordpress-mysql.yaml
└── utils
    ├── exec
    │   └── exec.go
    ├── os
    │   ├── file.go
    │   ├── readers.go
    │   └── writers.go
    ├── urlfile
    │   └── urlfile.go
    └── utils.go
```
