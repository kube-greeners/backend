<div id="top"></div>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<br />
<div align="center">

<h3 align="center">kube-green Dashboard Backend</h3>

  <p align="center">
    Making the kube-green savings stand out
    <br />
    <br />
    <a href="https://github.com/kube-greeners/backend/issues">Report Bug</a>
    ·
    <a href="https://github.com/kube-greeners/backend/issues">Request Feature</a>
  </p>
</div>


<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#used-technologies-and-libraries">Used Technologies and Libraries</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#local-installation-and-running">Local Installation and Running</a></li>
        <li><a href="#installation-to-a-kubernetes-cluster">Installation to a Kubernetes Cluster</a></li>
      </ul>
    </li>
    <li>
      <a href="#code-structure">Code Structure</a>
      <ul>
        <li><a href="#directories">Directories</a></li>
        <li><a href="#logic">Logic</a></li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#people">People</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

## About The Project

This project, as a part of the Dashboard for kube-green project intends to provide the backend API for executing specific Prometheus queries and returning the results to the frontend counterpart.
By itself this project is not that much of use and it would be better to combine it with the frontend to deploy.
This project is developed as a part of Distributed Software Development course taught at Politecnico di Milano and Mälardalen University in academic year 2021-2022, with the proposal provided by Mia Platform.
This project is built with the expectancy of accompanying [kube-green](https://github.com/kube-green/kube-green). 

<p align="right">(<a href="#top">back to top</a>)</p>

### Used technologies and Libraries

* [Go](https://go.dev/)
* [Prometheus](https://prometheus.io/)
* [K8S](https://kubernetes.io/)
* [Go CORS Handler](https://github.com/rs/cors)
* [Docker](https://www.docker.com/)
* [kube-green](https://github.com/kube-green/kube-green)
##### README is based on [Best-README-Template](https://github.com/othneildrew/Best-README-Template)
<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

This project assumes the availability of Prometheus, deployed and accessible from the application's host.
Furthermore, it assumes that certain Kubernetes and kube-green metrics are available.

To do this, directly provisioning Prometheus on top of Kubernetes using the Prometheus stack from Helm, and deploying kube-green on top will suffice.
Nevertheless, any preconfigured cluster with the [Kubernetes Node Exporter](https://github.com/prometheus/node_exporter) will also satisfy all the requirements.
To check what actual metrics are used, you can look at [queries.go](https://github.com/kube-greeners/backend/blob/dev/internal/queries.go).

### Prerequisites

To run this project on your system, you need:
* Go 1.17
* Prometheus endpoint URL accessible, which satisfies all the requirements given above.
* An empty port to serve the endpoint

### Local Installation and Running

1. Connect to the Prometheus endpoint if necessary, something like:
    ```sh
    kubectl port-forward <prometheus_pod> 9090
    ```
2. Clone the repo
   ```sh
   git clone https://github.com/kube-greeners/backend.git
   ```
3. Get Go packages
   ```sh
   go mod download
   ```
4. Export Prometheus and serve address
   ```sh
   export PROMETHEUS_URL="http://localhost:9090"
   export SERVE_ADDRESS=":8080"
   ```
5. Run
    ```sh
   go run ./cmd/server/main.go
    ```

<p align="right">(<a href="#top">back to top</a>)</p>

## Installation to a Kubernetes cluster

You can create a Deployment with the image `ghcr.io/kube-greeners/backend/backend`. An example deployment can be found below:
<details>

```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-green-backend
spec:
  replicas: 1
  selector:
    matchLabels:
      name: "kube-green-backend"
  template:
    metadata:
      labels:
        name: "kube-green-backend"
    spec:
      containers:
        - name: application
          image: ghcr.io/kube-greeners/backend/backend:latest
          env:
            - name: PROMETHEUS_URL
              value: http://prometheus-stack-kube-prom-prometheus.monitoring.svc.cluster.local:9090
            - name: SERVE_ADDRESS
              value: 0.0.0.0:8080
          ports:
            - containerPort: 8080
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: "/static"
              name: static
              readOnly: true
      volumes:
        - name: static
          emptyDir: {}

---
apiVersion: v1
kind: Service
metadata:
  name: kube-green-backend-svc
spec:
  selector:
    name: kube-green-backend
  type: ClusterIP
  ports:
    - port: 80
      targetPort: 8080
      name: http

```
</details>
<p align="right">(<a href="#top">back to top</a>)</p>

## Code structure

This repository consists of the main code for the backend, and some patches that we have performed on top to integrate with our workflow.
The Go code is ready to be executed on any Kubernetes cluster with kube-green, while the CI/CD code and helper scripts are run assuming our Kubernetes cluster exists.
The specific code to deploy to and link a specific cluster will probably be removed before making this project available to public.

### Directories
* `.github/`: Automation using Github Actions, configurations in here allow us to automatically test and deploy any change in the backend code to their appropriate branch.
* `cmd/`: Contains the bootstrapping code for the main code
* `internal/`: The main directory
  * `prometheus.go`: Code related to execution of Prometheus queries, abstractions to bootstrap the initialization and execution of queries to the Prometheus server.
  * `queries.go`: Go file containing the queries to be executed on the Prometheus server.
  * `server.go`: Code related to creating the HTTP server to serve different queries provided in `queries.go` file. It also can serve the `static` directory.
  * `*_test.go`: Unit tests
* `static/`: Directory to inject the frontend components to, the file inside exists as a stub for an actual application
* `connect.mjs`: A zx script that allows provisioning and connecting to our cluster backend, it helps developers of both frontend and backend projects to easily turn-on and turn-off the test cluster we have and port-forward the cluster pods.
* `get_frontend.py`: A python script that allows our Kubernetes backend deployment to automatically attach all the frontend branches that are built by the CI of the frontend
* `k8s-template.yml`: A K8S template allowing the provisioning of backend and frontend on top, so that we can directly create backends for each branch on Github Actions.

### Logic

The application in itself hosts a series of HTTP endpoints, that allows the user to query different resources to allow them observe their savings and decide on what to improve.
The main concern while building the application was to provide a backend for frontend, so the queries on the backend has evolved in parallel with the needs of the frontend project.
The detailed API can be seen by looking at the Swagger Documentation (WIP).
All the endpoints receive a start and end timestamp, to define the query boundaries.
The step parameter is calculated so that the endpoint will return at most 20 results to have a good visualization of data on the frontend.
Also, there is the optional namespace parameter allowing the frontend to filter the query results by namespace, which falls back to all namespaces as default.

<p align="right">(<a href="#top">back to top</a>)</p>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>

## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>

## People
#### Built by:
* Ozan İncesulu [@ozyinc](https://github.com/ozyinc)
* Alban Delishi [@albandelishi](https://github.com/albandelishi)
* Zoé Pesleux [@zoepj](https://github.com/zoepj)
* Redion Lila [@rlila97](https://github.com/rlila97)
#### [kube-green frontend](https://github.com/kube-greeners/frontend/) by:
* Boris Grunwald [@jikol1906](https://github.com/jikol1906)
* Ragnhild Kleiven [@RagnhildK](https://github.com/RagnhildK)
* Hanna Torjusen [@hanntorj](https://github.com/hanntorj)
* Marija Popovic [@marijapopovic28](https://github.com/marijapopovic28)
* Amila Cizmic [@amilacizmic](https://github.com/amilacizmic)
#### Acknowledgements:
We would like to thank:
* Malvina Latifaj and Samuele Giussani for their assistance and feedback through this project and for sharing this journey with us.
* Davide Bianchi for kube-green and valuable guidance for integration of this project
* Francesca Carta for feedback regarding each step of the product built.

<p align="right">(<a href="#top">back to top</a>)</p>

[contributors-shield]: https://img.shields.io/github/contributors/kube-greeners/backend.svg?style=for-the-badge
[contributors-url]: https://github.com/kube-greeners/backend/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/kube-greeners/backend.svg?style=for-the-badge
[forks-url]: https://github.com/kube-greeners/backend/network/members
[stars-shield]: https://img.shields.io/github/stars/kube-greeners/backend.svg?style=for-the-badge
[stars-url]: https://github.com/kube-greeners/backend/stargazers
[issues-shield]: https://img.shields.io/github/issues/kube-greeners/backend.svg?style=for-the-badge
[issues-url]: https://github.com/kube-greeners/backend/issues
[license-shield]: https://img.shields.io/github/license/kube-greeners/backend.svg?style=for-the-badge
[license-url]: https://github.com/kube-greeners/backend/blob/dev/LICENSE
