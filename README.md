<a name="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]


<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/vijethph/go-apps">
    <img src="https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg" width="100" height="100" >
  </a>

<h3 align="center">Go Apps</h3>

  <p align="center">
    Collection of things built following Go and Kubernetes tutorials
    <br />
    <br />
    <a href="https://github.com/vijethph/go-apps/issues">Report Bug</a>
    ·
    <a href="https://github.com/vijethph/go-apps/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#deploy-using-minikube">Deploy using minikube</a></li>
        <li><a href="#deploy-using-skaffold">Deploy using Skaffold</a></li>
        <li><a href="#teardown-for-locally-created-resources">Teardown for locally created resources</a></li>
        <li><a href="#deploy-to-gke">Deploy to GKE</a></li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

Collection of things built following Go and Kubernetes tutorials


### Built With

* [![Go][Go]][go-url]
* [![Docker][Docker]][docker-url]
* [![Kubernetes][Kubernetes]][kubernetes-url]


<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

The repository has multiple branches for `fibo-k8s` app based on deployment methods. 

- `main` branch contains files that can be run using Docker Compose (`docker-compose up`)
- `minikube-skaffold` branch contains files that can be run using minikube or Skaffold
- `eks-aks` branch contains files that can be deployed to AWS Elastic Kubernetes Service (EKS) or Azure Kubernetes Service (AKS)
- `gke-with-nginx` branch contains files that can be deployed to Google Kubernetes Engine (GKE)

### Deploy using minikube

1. Install kubectl and minikube. Then run these commands in order:
  ```bash
  minikube start

  minikube addons enable ingress

  kubectl cluster-info
  ```
2. Now, switch to minikube's docker daemon, as minikube runs in VM, and cannot use local images on its own. So, all docker images need to be built within minikube's docker daemon. Then, deploy the kubernetes cluster
  ```bash
  eval $(minikube docker-env)

  # docker build all images after above command

  kubectl create secret generic pgpassword --from-literal PG_PASSWORD=Test@123

  kubectl apply -f fibo-k8s

  kubectl get pods --watch

  minikube ip
  ```
3. Open the IP address provided by minikube in a browser to access the application. 

### Deploy using Skaffold

1. Follow the steps 1 and 2 from minikube deployment with certain changes:
  ```bash
  # instead of "kubectl apply -f k8s", run the following command
  
  skaffold dev # or skaffold run

  minikube ip
  ```
2. Open the IP address provided by minikube in a browser to access the application. 

### Teardown for locally created resources
```bash
kubectl delete -f fibo-k8s

skaffold delete

minikube stop

minikube delete --all

# switch back to normal docker daemon
eval $(minikube docker-env -u)
```

### Deploy to GKE

1. Enable Artifact Registry API and create a repository named `testrepo`
2. Enable Kubernetes Engine API. Create GKE cluster with necessary resources (which takes up to 15-20 minutes)
  ```
  gcloud auth login

  gcloud beta container clusters create "test-cluster" --no-enable-basic-auth --cluster-version "1.27.3-gke.100" --release-channel "regular" --machine-type "e2-medium" --image-type "COS_CONTAINERD" --disk-type "pd-standard" --disk-size "50" --metadata disable-legacy-endpoints=true --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes "3" --logging=SYSTEM,WORKLOAD --monitoring=SYSTEM --enable-ip-alias --no-enable-intra-node-visibility --default-max-pods-per-node "110" --security-posture=standard --workload-vulnerability-scanning=disabled --no-enable-master-authorized-networks --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --binauthz-evaluation-mode=DISABLED --no-enable-managed-prometheus --enable-shielded-nodes --node-locations "us-east1-d" --location us-east1
  ```
3. Install kubectl and Helm if GCP Cloud Shell is used to create resources
  ```
  gcloud components install kubectl

  curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
  
  chmod 700 get_helm.sh

  ./get_helm.sh
  ```
4. Build and push Docker images in `fibo-client`, `fibo-server` and `fibo-worker` folders to the Artifact Registry repository. Update the respective deployment files in `fibo-k8s` folder to reference the published images
5. Create a secret using kubectl to store database password. Install ingress-nginx controller for Network Load Balancer (NLB)
  ```
  gcloud container clusters get-credentials test-cluster --zone us-east-1

  kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value account)

  kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/aws/deploy.yaml
  ```
6. Create a secret using kubectl to store database password. Deploy the app to GKE
  ```
  kubectl create secret generic pgpassword --from-literal PG_PASSWORD=Test@123

  kubectl apply -f fibo-k8s
  ```
7. Get the external IP address of service, and open it in a new browser window to use the deployed application  
  ```
  kubectl get services --namespace=ingress-nginx
  ```
8. [Optional] Delete all the created resources by running the following commands
  ```
  kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/aws/deploy.yaml

  kubectl delete -f fibo-k8s

  gcloud container clusters delete test-cluster --location us-east1
  ```

**Alternative Method**: configure necessary credentials in GitHub repository secrets, complete steps 1 and 2 above and run the GitHub Actions workflow defined in this repository.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Vijeth P H - [@vijethph](https://github.com/vijethph)

Project Link: [https://github.com/vijethph/go-apps](https://github.com/vijethph/go-apps)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Official Go Documentation](https://go.dev/doc/)
* [Go Playground](https://go.dev/play/)
* [Go Standard Library Packages](https://pkg.go.dev/std)
* [Go By Example Tutorials](https://gobyexample.com/)
* [Go: The Complete Developer's Guide (Golang) by Stephen Grider](https://www.udemy.com/course/go-the-complete-developers-guide/)
* [Best-README-Template](https://github.com/othneildrew/Best-README-Template)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/vijethph/go-apps.svg?style=flat-square
[contributors-url]: https://github.com/vijethph/go-apps/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/vijethph/go-apps.svg?style=flat-square
[forks-url]: https://github.com/vijethph/go-apps/network/members
[stars-shield]: https://img.shields.io/github/stars/vijethph/go-apps.svg?style=flat-square
[stars-url]: https://github.com/vijethph/go-apps/stargazers
[issues-shield]: https://img.shields.io/github/issues/vijethph/go-apps.svg?style=flat-square
[issues-url]: https://github.com/vijethph/go-apps/issues
[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[go-url]: https://go.dev/doc/
[Docker]: https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white
[docker-url]: https://docs.docker.com/
[Kubernetes]: https://img.shields.io/badge/kubernetes-326ce5.svg?&style=for-the-badge&logo=kubernetes&logoColor=white
[kubernetes-url]: https://kubernetes.io/docs/home/