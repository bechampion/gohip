FROM ubuntu:18.04
RUN apt-get update && apt-get -y install gir1.2-gtk-3.0 gir1.2-webkit2-4.0 software-properties-common sudo && add-apt-repository ppa:yuezk/globalprotect-openconnect
RUN apt-get update && apt-get -y install globalprotect-openconnect
RUN useradd vpn
