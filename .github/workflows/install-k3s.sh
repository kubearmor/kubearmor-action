echo "RUNTIME="$RUNTIME

if [ "$RUNTIME" == "crio" ]; then
    ./install/self-managed-k8s/crio/install_crio.sh
fi

./install/k3s/install_k3s.sh