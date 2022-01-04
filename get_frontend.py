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


def is_branch_name(name_: str) -> bool:
    return name_.startswith("kg-") or name_ == "dev" or name_ == "master"


for configmap in configmaps:
    cm_name = configmap["metadata"]["name"]
    if is_branch_name(cm_name):
        directories[cm_name] = configmap["binaryData"]

os.mkdir("temp_")
for name, binaryData in directories.items():
    with open(f"temp_/{name}.tar.gz", "wb") as archive:
        archive.write(base64.decodebytes(bytes(binaryData['build.tar.gz'], "utf-8")))
    shutil.unpack_archive(f"temp_/{name}.tar.gz", f"temp_/{name}")
    shutil.copytree(f"temp_/{name}/build", f"static/{name}")

shutil.rmtree("temp_")
