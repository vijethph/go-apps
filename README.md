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
        <li><a href="#deploy-to-eks">Deploy to EKS</a></li>
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
- `eks-with-nlb` branch contains files that can be deployed to AWS Elastic Kubernetes Service (EKS)

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

### Deploy to EKS

Run the following commands with AWS CLI (using required credentials) or within AWS CloudShell
1. Install kubectl (v1.23.6), eksctl and helm.
2. Create EKS cluster with necessary resources (which takes up to 15-20 minutes)
  ```
  eksctl create cluster --name test-cluster --node-type t2.small --nodes 3 --nodes-min 3 --nodes-max 6 --region us-east-1
  ```
3. Create OIDC provider for the cluster, IAM service account and Amazon EBS CSI driver for the Persistent Volume Claim
  <code>
  eksctl create iamserviceaccount --region us-east-1 --name ebs-csi-controller-sa --namespace kube-system --cluster test-cluster --attach-policy-arn arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy --approve --role-only --role-name AmazonEKS_EBS_CSI_DriverRole

  eksctl create addon --name aws-ebs-csi-driver --cluster test-cluster --service-account-role-arn arn:aws:iam::$(aws sts get-caller-identity --query Account --output text):role/AmazonEKS_EBS_CSI_DriverRole --force
  </code>

4. Build and push Docker images in `fibo-client`, `fibo-server` and `fibo-worker` folders to a pre-created ECR repository. Update the respective deployment files in `fibo-k8s` folder to reference the published images
5. Create a secret using kubectl to store database password. Install ingress-nginx controller for Network Load Balancer (NLB) and finally deploy the app to EKS
  ```
  kubectl create secret generic pgpassword --from-literal PG_PASSWORD=Test@123

  kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/aws/deploy.yaml

  kubectl apply -f fibo-k8s
  ```
6. Get the DNS name of NLB and open it in a new browser window to use the deployed application
  ```
  kubectl get services --namespace=ingress-nginx
  ```
6. [Optional] Delete all the created resources by running the following commands
  ```
  kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/aws/deploy.yaml

  kubectl delete -f fibo-k8s

  eksctl delete cluster --name test-cluster --region us-east-1

  eksctl delete iamserviceaccount --name ebs-csi-controller-sa --region us-east-1
  ```

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