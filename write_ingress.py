import json
import subprocess

a = subprocess.run("kubectl get svc -n backend -o name".split(" "), capture_output=True).stdout
branches = [res[8:-4] for res in a.decode("utf-8").split("\n") if res]
path_specs = []
for branch in branches:
    path_specs.append({
        "host": f"{branch}.dsd.ozyinc.com",
        "http": {
            "paths": [{
                "path": "/*",
                "pathType": "ImplementationSpecific",
                "backend": {
                    "service": {
                        "name": f"{branch}-svc",
                        "port": {
                            "number": 80
                        }
                    }
                }
            }]
        }
    })
spec = {
    "apiVersion": "networking.k8s.io/v1",
    "kind": "Ingress",
    "metadata": {
        "name": "test",
        "namespace": "backend",
        "annotations": {
            "kubernetes.io/ingress.class": "gce",
            "kubernetes.io/ingress.global-static-ip-name": "dsd-ip",
        }
    },
    "spec": {
        "rules": path_specs,
    }
}
json.dump(spec, open("ingress.json", "w"))
