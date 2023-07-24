#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright 2023 Authors of KubeArmor


echo "RUNTIME="$RUNTIME

if [ "$RUNTIME" == "crio" ]; then
    ./install/self-managed-k8s/crio/install_crio.sh
fi

./install/k3s/install_k3s.sh