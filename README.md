<p>This is a backend project I made after I learnt how to use Go language and Echo framework. This application is made on microservices architecture
. This project consists of two services. One is the steganography encoding service. The other is the message encryption service. These app image are also available in dockerhub in this <a href="https://hub.docker.com/r/rixelanya/crypstego">link</a></p>

<h4>How Does This Work</h4>
These apps were compiled and then made into Docker image. Then these images are containerized and orchestrated using Kubernetes. My Kubernetes configuration yaml configure them so that two pods or replica of each app is running. And also Kubernetes expose these services at port 30007 of the Kubernetes Controller Machine.

<h5>How To use </h5>


First clone this repoitory to your machine and then make sure you have kubernetes and docker on your machine <br></br>
Then apply the following command <br></br>
kubectl apply -f crypstegokube.yaml


After applying this go to http://localhost:30007/encode.html
