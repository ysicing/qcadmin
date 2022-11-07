#!/usr/bin/env sh

helm repo add install https://hub.qucheng.com/chartrepo/stable

helm repo update

kubectl create ns cne-system

export APP_DOMAIN=${APP_DOMAIN:-k3s.local}
export APP_TOKEN=${APP_TOKEN:-XZdrjxhAhq5pDjpEU3kR4djsvJ3rfj0M}
export TOP_DOMAIN=${APP_DOMAIN#*.}

kubectl apply -f  https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml

cat > /tmp/qucheng.yaml <<EOF
cloud:
  defaultChannel: stable
env:
  APP_DOMAIN: ${APP_DOMAIN}
  CNE_API_TOKEN: ${APP_TOKEN}
ingress:
  host: ${APP_DOMAIN}
EOF

helm upgrade -i ingress install/nginx-ingress-controller -n cne-system
helm upgrade -i cne-operator install/cne-operator -n cne-system
helm upgrade -i qucheng install/qucheng -f /tmp/qucheng.yaml -n cne-system

[ -d "/qcli/root/.kube" ] || mkdir -pv /qcli/root/.kube
[ -d "/qcli/root/.qc/config" ] || mkdir -pv /qcli/root/.qc/config

cp -a /qcadmin_linux_amd64 /qcli/qbin/q
cp -a /qcadmin_linux_amd64 /qcli/qbin/qcadmin
cp -a /usr/local/bin/kubectl /qcli/qbin/kubectl
cp -a /usr/local/bin/helm /qcli/qbin/helm

cat > /qcli/root/.qc/config/cluster.yaml <<EOF
api_token: ${APP_TOKEN}
cluster_id: ""
console-password: pass4Quickon
db: sqlite
domain: ${APP_DOMAIN}
init_node: ${APP_NODE_IP}
master:
- host: ${APP_NODE_IP}
  name: ${APP_NODE_IP}
token: ""
worker: null
EOF

[ -f "/qcli/k3syaml/k3s.yaml" ] && cp -a /qcli/k3syaml/k3s.yaml /qcli/root/.kube/config