grep -rhE ^deb /etc/apt/sources.list* | grep "cloud-sdk"
sudo apt-get update
sudo apt-get install -y kubectl
kubectl version --client