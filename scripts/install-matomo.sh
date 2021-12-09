#!/usr/bin/env bash

# run as root

apt-get update -y
apt-get install curl wget vim git unzip socat gnupg2 -y
apt-get install nginx mariadb-server php7.4 php7.4-cli php7.4-fpm php7.4-common php7.4-curl php7.4-gd php7.4-xml php7.4-mbstring php7.4-mysql -y

mysql -e "CREATE DATABASE matomodb"
mysql -e "GRANT ALL ON matomodb.* TO 'matomo' IDENTIFIED BY 'password'"
mysql -e "FLUSH PRIVILEGES"
cd /var/www/html/
wget https://builds.matomo.org/matomo-latest.zip
unzip matomo-latest.zip
chown -R www-data:www-data /var/www/html/matomo
rm matomo-latest.zip

cat  > /etc/nginx/sites-available/matomo.conf << EOF
server { # remove this if you don't want Matomo to be reachable from IPv6
    listen 80;
    server_name matomo.geogame.xyz;
    add_header Referrer-Policy origin always; # make sure outgoing links don't show the URL to the Matomo instance
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    root /var/www/html/matomo/;
    index index.php;
    
    ## only allow accessing the following php files
    location ~ ^/(index|matomo|piwik|js/index|plugins/HeatmapSessionRecording/configs)\.php$ {
        include snippets/fastcgi-php.conf; # if your Nginx setup doesn't come with a default fastcgi-php config, you can fetch it from https://github.com/nginx/nginx/blob/master/conf/fastcgi.conf
        #try_files $fastcgi_script_name =404; # protects against CVE-2019-11043. If this line is already included in your snippets/fastcgi-php.conf you can comment it here.
        fastcgi_param HTTP_PROXY ""; # prohibit httpoxy: https://httpoxy.org/
    	fastcgi_pass unix:/var/run/php/php7.4-fpm.sock; 
        #fastcgi_pass 127.0.0.1:9000; # uncomment if you are using PHP via TCP sockets (e.g. Docker container)
    }

    ## deny access to all other .php files
    location ~* ^.+\.php$ {
        deny all;
        return 403;
    }

    ## serve all other files normally
    location / {
        try_files $uri $uri/ =404;
    }

    ## disable all access to the following directories
    location ~ ^/(config|tmp|core|lang) {
        deny all;
        return 403; # replace with 404 to not show these directories exist
    }

    location ~ /\.ht {
        deny  all;
        return 403;
    }

    location ~ js/container_.*_preview\.js$ {
        expires off;
        add_header Cache-Control 'private, no-cache, no-store';
    }

    location ~ \.(gif|ico|jpg|png|svg|js|css|htm|html|mp3|mp4|wav|ogg|avi|ttf|eot|woff|woff2|json)$ {
        allow all;
        ## Cache images,CSS,JS and webfonts for an hour
        ## Increasing the duration may improve the load-time, but may cause old files to show after an Matomo upgrade
        expires 1h;
        add_header Pragma public;
        add_header Cache-Control "public";
    }

    location ~ ^/(libs|vendor|plugins|misc|node_modules) {
        deny all;
        return 403;
    }

    ## properly display textfiles in root directory
    location ~/(.*\.md|LEGALNOTICE|LICENSE) {
        default_type text/plain;
    }
}
EOF
ln -s /etc/nginx/sites-available/matomo.conf /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
