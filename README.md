git clone https://github.com/Srgkharkov/recordingmeet.git

docker run -it --rm \
  -v /root/recordingmeet/nginx/certs:/etc/letsencrypt \
  -v /root/recordingmeet/nginx:/var/lib/letsencrypt \
  -p 80:80 \
  certbot/certbot certonly --standalone --standalone -d rec.srgkharkov.ru

  docker run -p 8080:8080 -e JWT_SECRET_KEY=my_super_secret_key rec