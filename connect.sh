
# Grafana
kubectl -n monitoring port-forward $(kubectl get pod -n monitoring -l "app.kubernetes.io/name=grafana" -o jsonpath="{.items[0].metadata.name}") 3000 >/dev/null &
echo "Grafana: http://localhost:3000"
echo "Username: admin"
echo "Password: test"
echo ""

# Prometheus
kubectl -n monitoring port-forward $(kubectl get pod -n monitoring -l "app.kubernetes.io/name=prometheus" -o jsonpath="{.items[0].metadata.name}") 9090 >/dev/null&
echo "Prometheus: http://localhost:9090"
echo ""

# Alert Manager
kubectl -n monitoring port-forward $(kubectl get pod -n monitoring -l "app.kubernetes.io/name=alertmanager" -o jsonpath="{.items[0].metadata.name}") 9093 >/dev/null&
echo "Alert Manager: http://localhost:9093"

export PROMETHEUS_URL="http://localhost:9090"
wait
