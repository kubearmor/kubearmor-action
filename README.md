# kubearmor-action
A Github Action used to identify changes in the application posture. Such as what new processes are being spawned and what new file accesses are being made.

## How To Use
```yaml
name: example
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
jobs:
  example_job:
    runs-on: ubuntu-latest
    name: A job to test kubearmor-action
    steps:
      - name: kubearmor-action test
        uses: kubearmor/kubearmor-action@main
        with:
          old-app: 'https://raw.githubusercontent.com/{user}/{repo}/main/examples/wordpress-mysql/wordpress-mysql-deployment.yaml'
          new-app: 'manifests/app.yaml'
          namespace: 'wordpress-mysql'
          app-name: 'wordpress'
      - uses: actions/download-artifact@v2
        with:
          name: ${{ steps.visualisation-report.outputs.visualisation-results-artifact }}
          path: images
      - uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./images
          keep_files: true
      - name: comment PR
        run: |
          gh pr comment  ${{ github.event.number }} -b"![system_graph](https://raw.githubusercontent.com/${{ github.repository }}/gh-pages/${{ steps.visualisation-report.outputs.sys-visualisation-image }})"
          gh pr comment  ${{ github.event.number }} -b"![network_graph](https://raw.githubusercontent.com/${{ github.repository }}/gh-pages/${{ steps.visualisation-report.outputs.network-visualisation-image }})"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
## Architecture Overview
![Alt text](doc/pics/kubearmor-action-Architecture.drawio.png)