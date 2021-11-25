#!/usr/bin/env zx

$.verbose = false;

const clearLastLine = () => {
    process.stdout.moveCursor(0, -1) // up one line
    process.stdout.clearLine(1) // from cursor to end
    // https://stackoverflow.com/questions/32938213/is-there-a-way-to-erase-the-last-line-of-output/46510907#46510907
}

const capitalize = s => s && s[0].toUpperCase() + s.slice(1) // https://stackoverflow.com/questions/1026069/how-do-i-make-the-first-letter-of-a-string-uppercase-in-javascript

const scale = async (count) => {
    await $`gcloud container clusters resize test-cluster --project kube-greeners --zone europe-west1-b --node-pool optimized-pool --num-nodes ${count} --async --quiet`;
}

const waitForResources = async (count) => {
    console.log("");
    let nodes;
    do {
        nodes = JSON.parse((await $`kubectl get nodes -o json`)).items;
        // Nodes are being created
        if (nodes.length >= count) {
            break;
        }
        clearLastLine();
        console.log(`Waiting to provision enough nodes: ${nodes.length}/${count}`);
        await sleep(2000);

    } while (nodes.length < count)

    let readyNodes;
    do {
        nodes = JSON.parse((await $`kubectl get nodes -o json`)).items;
        // The nodes that are ready to schedule
        readyNodes = nodes.filter((node) => {
            return node.status.conditions
                .filter((condition) => condition.reason === "KubeletReady" && condition.status === "True")
                .length > 0;
        });
        if (readyNodes.length === nodes.length) {
            break;
        }
        clearLastLine();
        console.log(`Waiting for nodes: ${readyNodes.length} / ${nodes.length}`);
        await sleep(2000);
    } while (readyNodes.length !== nodes.length);


    // Monitoring and backend deployments are ready
    let workloads;
    let readyWorkloads;
    do {
        workloads = (await Promise.all([
            $`kubectl -n monitoring get deployments -o json`,
            $`kubectl -n monitoring get statefulsets -o json`,
            $`kubectl -n backend get deployments -o json`,
        ])).map(JSON.parse).map(resource => resource.items).flat();
        readyWorkloads = workloads
            .filter((workload) => workload.status.replicas === workload.status.readyReplicas);
        if (readyWorkloads.length === workloads.length) {
            break;
        }
        clearLastLine();
        console.log(`Waiting for workloads to be active ${readyWorkloads.length} / ${workloads.length}`);
        await sleep(2000);
    } while (readyWorkloads.length !== workloads.length);

    // All done, let's go
    clearLastLine();
    console.log("All resources are ready.");
}
// Check if nodes are available
const nodeCheck = async () => {
    let nodeCount = JSON.parse(await $`kubectl get nodes -o json`).items.length

    if (nodeCount === 0) {
        console.log(chalk.red("There are no nodes ready"))
        const requestedNodeCount = Number.parseInt(await question("How many nodes shall I provision? [3] ")) || 3;
        console.log("Ok, this will take time, grab coffee :)");
        await scale(requestedNodeCount);
        await waitForResources(requestedNodeCount);
    } else {
        await waitForResources(1); // We just want some nodes if they are already provisioned
    }
}

const processPromises = [];
// Monitoring tooling connection
const tools = [{
    name: "grafana",
    port: 9092,
    extraInfo: `Username: admin\nPassword: test`,
}, {
    name: "prometheus",
    port: 9090,
}, {
    name: "alertmanager",
    port: 9091,
}];
const provisionTool = async tool => {

    const url = `http://localhost:${tool.port}`
    process.env[tool.name.toUpperCase() + "_URL"] = url;
    const podName = await $`kubectl get pod -n monitoring -l "app.kubernetes.io/name=${tool.name}" -o jsonpath="{.items[0].metadata.name}"`;
    const portForwardPromise = $`kubectl -n monitoring port-forward ${podName} ${tool.port}`;
    processPromises.push(portForwardPromise);
    console.log(chalk.green(`${capitalize(tool.name)}: ${url}`));
    if (tool.extraInfo) {
        console.log(chalk.green(tool.extraInfo));
    }
    console.log("");
    await portForwardPromise;
}

// Backend connection
const provisionBackend = async (backendVersion) => {

    const backendPort = "8080";
    if (backendVersion === "local") {
        process.env["SERVE_ADDRESS"] = ":" + backendPort;
        const goPromise = $`go run cmd/server/main.go`;
        processPromises.push(goPromise);
        console.log(chalk.blue(`Backend: http://localhost:${backendPort} (from local code)`));
        await goPromise;
    } else {
        if (await $`kubectl -n backend describe deployment ${backendVersion}`.exitCode !== 0) {
            console.error(chalk.red(`Backend not found: ${backendVersion}`));
            process.exit(-1);
        }
        const podName = await $`kubectl get pod -n backend -l "name=${backendVersion}" -o jsonpath="{.items[0].metadata.name}"`;
        const backendPortForwardPromise = $`kubectl -n backend port-forward ${podName} ${backendPort}`;
        processPromises.push(backendPortForwardPromise);
        console.log(chalk.blue(`Backend: http://localhost:${backendPort} (for branch: ${backendVersion})`));
        await backendPortForwardPromise;
    }
}

// Ask if the user  wants to shut down cluster
process.on('SIGINT', async () => {
    processPromises.map(process => process.kill());
    console.clear();
    const downscale = await question("Shall I downscale the cluster? [Y/n] ") || "y";
    if (downscale === "y" || downscale === "Y") {
        console.log("Ok, triggered downscale, shouldn't take long");
        await scale(0);
    } else {
        console.log("Ok, cluster is as it is");
    }
    console.log("Bye!");
    process.exit(0);
});


// Main
const asyncPromises = [];
await nodeCheck();
asyncPromises.push(...tools.map(provisionTool));
if (argv.backend) {
    const backendVersion = argv.backend.toLowerCase();
    asyncPromises.push(provisionBackend(backendVersion));
} else {
    console.log(chalk.blue("Not connecting to a backend"));
}
try {
    await Promise.all(asyncPromises);
} catch (e) {
    processPromises.map(process => process.kill());
}


