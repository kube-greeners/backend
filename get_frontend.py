import os
import json
import shutil
import base64

# Point to the internal API server hostname
APISERVER = "https://kubernetes.default.svc"

# Path to ServiceAccount token
SERVICEACCOUNT = "/var/run/secrets/kubernetes.io/serviceaccount"

# Read the ServiceAccount bearer token
TOKEN_PATH = f"{SERVICEACCOUNT}/token"

with open(TOKEN_PATH, "r") as f:
    TOKEN = f.read()

# Reference the internal certificate authority (CA)
CACERT_PATH = f"{SERVICEACCOUNT}/ca.crt"

# Explore the API with TOKEN
command = f'curl --cacert {CACERT_PATH} --header "Authorization: Bearer {TOKEN}" -X GET {APISERVER}/api/v1/namespaces/frontend/configmaps'
configmaps_resp = os.popen(command).read()
configmaps = json.loads(configmaps_resp)["items"]
directories = {}
print(configmaps.keys())
for configmap in configmaps:
    if configmap["metadata"]["name"].startswith("kg-"):
        directories[configmap["metadata"]["name"]] = configmap["binaryData"]

for name, binaryData in directories.items():
    with open(f"{name}.tar.gz", "wb") as archive:
        archive.write(base64.decodebytes(bytes(binaryData['build.tar.gz'], "utf-8")))
    shutil.unpack_archive(f"{name}.tar.gz", f"{name}")
    shutil.copytree(f"{name}/build", f"static/{name}")
    shutil.rmtree(f"{name}")
    os.remove(f"{name}.tar.gz")
print(os.popen("ls -al ."))
print(os.popen("ls -al ./static"))
