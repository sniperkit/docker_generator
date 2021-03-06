// This is simple Dockerfile generator
package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

)

type ParsingYaml struct {
	From     string      `yaml:"FROM"`               //distrib
	Composer string      `yaml:"composer,omitempty"` //install compose YES or NO
	Nginx    string      `yaml:"nginx,omitempty"`    //install nginx YES or NO
	PhpExt   interface{} `yaml:"php_modules"`        // list php-modules
}

type Version struct {
	php         string
	distrib     string
	packageName string
}

//install all-software
func SoftInstallApk(nginxSw bool) string {
	software := "RUN apk add --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing gnu-libiconv && \\\n"
	software += "echo @testing http://nl.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories && \\\n"
	software += "echo @main http://mirror.yandex.ru/mirrors/alpine/edge/main >>  /etc/apk/repositories && \\\n"
	software += "echo @community http://mirror.yandex.ru/mirrors/alpine/edge/community >>  /etc/apk/repositories && \\\n"
	software += "echo /etc/apk/repositories && \\\n"
	software += "apk update && \\\n"
	software += "apk add --no-cache bash \\\n"
	if nginxSw {
		software += "nginx \\\n"
	}
	software += "wget \\\n"
	software += "supervisor \\\n"
	software += "curl \\\n"
	software += "libcurl \\\n"
	software += "git \\\n"
	software += "python \\\n"
	software += "python-dev \\\n"
	software += "py-pip \\\n"
	software += "augeas-dev \\\n"
	software += "openssl-dev \\\n"
	software += "ca-certificates \\\n"
	software += "dialog \\\n"
	software += "autoconf \\\n"
	software += "make \\\n"
	software += "gcc \\\n"
	software += "musl-dev \\\n"
	software += "linux-headers \\\n"
	software += "libmcrypt-dev \\\n"
	software += "libpng-dev \\\n"
	software += "icu-dev \\\n"
	software += "libpq \\\n"
	software += "libxslt-dev \\\n"
	software += "libffi-dev \\\n"
	software += "freetype-dev \\\n"
	software += "sqlite-dev \\\n"
	software += "bzip2-dev \\\n"
	software += "libmemcached-dev \\\n"
	software += "libjpeg-turbo-dev \\\n"
	software += "&& \\\n"
	return software
}

// php composer installation
func PhpComposerSetup() string {
	compose := "EXPECTED_COMPOSER_SIGNATURE=$(wget -q -O - https://composer.github.io/installer.sig) && \\\n"
	compose += "	php -r \"copy('https://getcomposer.org/installer', 'composer-setup.php');\" && \\\n"
	compose += "	php -r \"if (hash_file('SHA384', 'composer-setup.php') === '${EXPECTED_COMPOSER_SIGNATURE}') "
	compose += "{ echo 'Composer.phar Installer verified'; } else { echo 'Composer.phar Installer corrupt'; "
	compose += "unlink('composer-setup.php'); } echo PHP_EOL;\" && \\\n"
	compose += "php composer-setup.php --install-dir=/usr/bin --filename=composer && \\\n"
	compose += "php -r \"unlink('composer-setup.php');\" && \\\n"
	return compose
}

// make - make install for php_modules
func StdConfAndMake() string {
	MakeAndConf := "phpize &&\\\n"
	MakeAndConf += "./configure && \\\n"
	MakeAndConf += "make && \\\n"
	MakeAndConf += "make install &&\\\n"
	return MakeAndConf
}

//install memcached
func InstallMemcached() string {
	memcach := "apk add --virtual .memcached-build-dependencies \\\n"
	memcach += "	libmemcached-dev \\\n"
	memcach += "	cyrus-sasl-dev && \\\n"
	memcach += "apk add --virtual .memcached-runtime-dependencies \\\n"
	memcach += "libmemcached &&\\\n"
	memcach += "git clone -o ${MEMCACHED_TAG} --depth 1 https://github.com/php-memcached-dev/php-memcached.git /tmp/php-memcached && \\\n"
	memcach += "cd /tmp/php-memcached &&\\\n"
	memcach += "phpize &&\\\n"
	memcach += "./configure \\\n"
	memcach += "    --disable-memcached-sasl \\\n"
	memcach += "    --enable-memcached-msgpack \\\n"
	memcach += "    --enable-memcached-json && \\\n"
	memcach += "make && \\\n"
	memcach += "make install && \\\n"
	memcach += "apk del .memcached-build-dependencies && \\\n"
	return memcach
}

//install msgpack
func InstallMsgpack() string {
	msgpack := "git clone -o ${MSGPACK_TAG} --depth 1 https://github.com/msgpack/msgpack-php.git /tmp/msgpack-php && \\\n"
	msgpack += "cd /tmp/msgpack-php && \\\n"
	msgpack += StdConfAndMake()
	return msgpack
}

//Install Imagick
func InstallImagick() string {
	imagick := "apk add --no-cache --virtual .imagick-build-dependencies \\\n"
	imagick += "  autoconf \\\n"
	imagick += "  g++ \\\n"
	imagick += "  gcc \\\n"
	imagick += "  git \\\n"
	imagick += "  imagemagick-dev \\\n"
	imagick += "  libtool \\\n"
	imagick += "  make \\\n"
	imagick += "  tar && \\\n"
	imagick += "apk add --virtual .imagick-runtime-dependencies \\\n"
	imagick += "  imagemagick &&\\\n"
	imagick += "git clone -o ${IMAGICK_TAG} --depth 1 https://github.com/mkoppanen/imagick.git /tmp/imagick &&\\\n"
	imagick += "cd /tmp/imagick && \\\n"
	imagick += StdConfAndMake()
	imagick += "echo \"extension=imagick.so\" > /usr/local/etc/php/conf.d/ext-imagick.ini && \\\n"
	imagick += "apk del .imagick-build-dependencies && \\\n"
	return imagick
}

//Install Xdebug
func InstallXdebug() string {
	xdebug := "git clone -o ${XDEBUG_TAG} --depth 1 https://github.com/xdebug/xdebug.git /tmp/xdebug &&\\\n"
	xdebug += "cd /tmp/xdebug && \\\n"
	xdebug += StdConfAndMake()
	return xdebug
}

//Install Redis
func InstallRedis() string {
	redis := "git clone -o ${REDIS_TAG} --depth 1 https://github.com/phpredis/phpredis.git /tmp/redis &&\\\n"
	redis += "cd /tmp/redis \\\n"
	redis += StdConfAndMake()
	return redis
}

// Setup modules from code(GIT)
func UnstandartModulesInstall(module string) (string, string) {
	if module == "memcached" {
		// version memcached
		arg := "ARG MEMCACHED_TAG=v3.0.4"
		return arg, InstallMemcached()
	}
	if module == "msgpack" {
		// version msgpack
		arg := "ARG MSGPACK_TAG=msgpack-2.0.2"
		return arg, InstallMsgpack()
	}
	if module == "imagick" {
		// version imagick
		arg := "ARG IMAGICK_TAG=\"3.4.2\""
		return arg, InstallImagick()
	}
	if module == "xdebug" {
		// version xdebug
		arg := "ARG XDEBUG_TAG=2.6.0"
		return arg, InstallXdebug()
	}
	if module == "redis" {
		// version redis
		arg := "ARG REDIS_TAG=3.1.6"
		return arg, InstallRedis()
	}
	return "", ""
}

//Create Dockerfile from alpine-based php-fpm
func Alpine(phpModules []interface{},
	ModulesNopecl []string,
	DockerModules []string,
	PhpVersion map[string]Version,
	maintainer string,
	confYaml ParsingYaml,
) string {
	var switcher bool
	ARG := "\n"
	ModulesLines := ""
	DockerPhpExtInstall := "docker-php-ext-install "
	GDconf := ""
	var arg, StrModule string
	for _, module := range phpModules {
		strModule := module.(string)
		switcher = true
		for _, nopecl := range ModulesNopecl {
			if strModule == nopecl {
				arg, StrModule = UnstandartModulesInstall(strModule)
				//generate script
				switcher = false
				ModulesLines += StrModule
				ARG += arg + "\n"
			}
		}
		if switcher {
			for _, dockerModule := range DockerModules {
				if strModule == dockerModule {
					DockerPhpExtInstall += strModule + " "
					if strModule == "gd" {
						GDconf += "docker-php-ext-configure gd \\\n"
						GDconf += "--with-gd \\\n"
						GDconf += "--with-freetype-dir=/usr/include/ \\\n"
						GDconf += "--with-png-dir=/usr/include/ \\\n"
						GDconf += "--with-jpeg-dir=/usr/include/ && \\\n"
					}
				}

			}
		}

	}
	nginxSW := false
	if confYaml.Nginx == "YES" {
		nginxSW = true
	}
	letsencrypt := "pip install -U pip && \\\npip install -U certbot && \\\nmkdir -p /etc/letsencrypt/webrootauth && \\\n"
	DockerPhpExtInstall += "&& \\\ndocker-php-source delete && \\\n"
	ARG += "\n"
	clean := "apk del gcc musl-dev linux-headers libffi-dev augeas-dev python-dev make autoconf \n"
	ENV := "ENV php_conf /usr/local/etc/php-fpm.conf\n"
	ENV += "ENV fpm_conf /usr/local/etc/php-fpm.d/www.conf\n"
	ENV += "ENV php_vars /usr/local/etc/php/conf.d/docker-vars.ini\n"
	ENV += "ENV LD PRELOAD /usr/lib/preloadable_libconv.so php\n"
	HEAD := "FROM " + PhpVersion[confYaml.From].packageName + "\n" + "LABEL maintainer = " + maintainer + "\n"

	Dockerfile := HEAD
	Dockerfile += ENV
	Dockerfile += ARG
	if PhpVersion[confYaml.From].distrib == "alpine" {
		Dockerfile += SoftInstallApk(nginxSW)
	}
	Dockerfile += GDconf
	Dockerfile += DockerPhpExtInstall
	if confYaml.Composer == "YES" {
		Dockerfile += PhpComposerSetup()
	}
	Dockerfile += ModulesLines
	Dockerfile += letsencrypt
	Dockerfile += clean
	ConfigList := "ADD supervisord.conf /etc/supervisor.conf"
	ConfigList += "ADD start.sh /start.sh"
	if nginxSW {
		ConfigList += "EXPOSE 443 80"
	} else {
		ConfigList += "EXPOSE 9000"
	}
	Dockerfile += ConfigList
	err := CreateSupervisord(nginxSW)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return Dockerfile
}

//generate supervisor config
func CreateSupervisord(nginxsw bool) error {
	supervisor := "[unix_http_server]\n"
	supervisor += "file=/dev/shm/supervisor.sock ; (the path to the socket file)\n"
	supervisor += "\n"
	supervisor += "[supervisord]\n"
	supervisor += "logfile=/tmp/supervisord.log ; (main log file;default $CWD/supervisord.log)\n"
	supervisor += "logfile_maxbytes=50MB ; (max main logfile bytes b4 rotation;default 50MB)\n"
	supervisor += "logfile_backups=10 ; (num of main logfile rotation backups;default 10)\n"
	supervisor += "loglevel=info ; (log level;default info; others: debug,warn,trace)\n"
	supervisor += "pidfile=/tmp/supervisord.pid ; (supervisord pidfile;default supervisord.pid)\n"
	supervisor += "nodaemon=false ; (start in foreground if true;default false)\n"
	supervisor += "minfds=1024 ; (min. avail startup file descriptors;default 1024)\n"
	supervisor += "minprocs=200 ; (min. avail process descriptors;default 200)\n"
	supervisor += "user=root ;\n"
	supervisor += "\n"
	supervisor += "; the below section must remain in the config file for RPC\n" +
		"; (supervisorctl/web interface) to work, additional interfaces may be\n" +
		"; added by defining them in separate rpcinterface: sections"
	supervisor += "[rpcinterface:supervisor]\n"
	supervisor += "supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface\n"
	supervisor += "\n[supervisorctl]\n"
	supervisor += "serverurl=unix:///dev/shm/supervisor.sock ; use a unix:// URL for a unix socket\n"
	//start php-fpm
	supervisor += "\n[program:php-fpm]\n"
	supervisor += "command = /usr/local/sbin/php-fpm --nodaemonize --fpm-config /usr/local/etc/php-fpm.d/www.conf\n"
	supervisor += "autostart=true\n"
	supervisor += "autorestart=true\n"
	supervisor += "priority=5\n"
	supervisor += "stdout_logfile=/dev/stdout\n"
	supervisor += "stdout_logfile_maxbytes=0\n"
	supervisor += "stderr_logfile=/dev/stderr\n"
	supervisor += "stderr_logfile_maxbytes=0\n"
	//start nginx
	if nginxsw {
		supervisor += "\n[program:nginx]\n"
		supervisor += "command=/usr/sbin/nginx -g \"daemon off; error_log /dev/stderr info;\"\n"
		supervisor += "autostart=true\n"
		supervisor += "autorestart=true\n"
		supervisor += "priority=10\n"
		supervisor += "stdout_events_enabled=true\n"
		supervisor += "stderr_events_enabled=true\n"
		supervisor += "stdout_logfile=/dev/stdout\n"
		supervisor += "stdout_logfile_maxbytes=0\n"
		supervisor += "stderr_logfile=/dev/stderr\n"
		supervisor += "stderr_logfile_maxbytes=0\n"
	}
	supervisor += "\n[include]\n"
	supervisor += "files = /etc/supervisor/conf.d/*.conf\n"
	ioutil.WriteFile("supervisord.conf", []byte(supervisor), 0755)
	return nil
}
//Create run script
func GenerateRunScript() error {
	//This function generate script.sh
	script := "adyn"
	script += GenerateCustomConfigurationPhp()
	return nil
}

//Configure custom configuration to php
func GenerateCustomConfigurationPhp() string {
	//This function generate custom config php.ini
	parameters := make(map[string]string)
	uploadMaxFilesize := "100M"
	postMaxSize := "100M"
	MemoryLimit := "128M"
	parameters["UploadMaxFilesize"] = uploadMaxFilesize
	parameters["PostMaxSize"] = postMaxSize
	parameters["MemoryLimit"] = MemoryLimit
	phpConfigure := "RUN " + "echo \"upload_max_filesize =  " + parameters["UploadMaxFilesize"] + "\" >> ${php_vars} &&\\\n"
	phpConfigure += "echo \"post_max_size =  " + parameters["PostMaxSize"] + "\" >> ${php_vars} &&\\\n"
	phpConfigure += "echo \"memory_limit =  " + parameters["MemoryLimit"] + "\" >> ${php_vars} &&\\\n"
	return phpConfigure
}

//Create Dockerfile from debian-based php-fpm
func Debian(phpModules []interface{},
	ModulesNopecl []string,
	DockerModules []string,
	PhpVersion map[string]Version,
	maintainer string,
	confYaml ParsingYaml,
) string {
	aptlib := make(map[string]string)
	aptlib["xsl"] = "libxslt-dev \\\n"
	aptlib["mcrypt"] = "libmcrypt \\\n"
	aptlib["memcached"] = "libmemcached-dev \\\nzlib1g-dev \\\n"
	aptlib["pdo_sqlite"] = "libsqlite3-dev \\\n"
	aptlib["intl"] = "libicu-dev \\\n"
	aptlib["bz2"] = "libbz2-dev \\\n"
	apt := "RUN apt-get update && apt-get install -y \\\n"
	DockerPhpExtInstall := "&&docker-php-ext-install "
	for _, modules := range phpModules {
		strmodule := modules.(string)
		for _, phpmodule := range DockerModules {
			if strmodule == phpmodule {
				DockerPhpExtInstall += strmodule + " "
				apt += aptlib[strmodule]
			}
		}
		for _, pecl := range ModulesNopecl {
			if strmodule == pecl {
				apt += aptlib[strmodule]
				//pecl
			}
		}


		}

	HEAD := "FROM " + PhpVersion[confYaml.From].packageName + "\n" + "LABEL maintainer = " + maintainer + "\n"
	Dockerfile := HEAD
	Dockerfile += apt + DockerPhpExtInstall
	return Dockerfile
}
func main() {
	conf, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//парсим YAML
	confYaml := ParsingYaml{}
	err = yaml.Unmarshal([]byte(conf), &confYaml)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	PhpVersion := make(map[string]Version)
	PhpVersion["7.1-alpine"] = Version{
		php:         "7.1",
		distrib:     "alpine",
		packageName: "php:7.1-fpm-alpine",
	}
	PhpVersion["7.2-alpine"] = Version{
		php:         "7.1",
		distrib:     "alpine",
		packageName: "php:7.2-fpm-alpine",
	}
	PhpVersion["7.1-fpm"] = Version{
		php:         "7.1",
		distrib:     "debian",
		packageName: "php:7.1-fpm",
	}
	PhpVersion["7.2-fpm"] = Version{
		php:         "7.2",
		distrib:     "debian",
		packageName: "php:7.2-fpm",
	}

	maintainer := "\"DockerFile generator by fp <alexwolk01@gmail.com>\" \n"
	phpModules, _ := confYaml.PhpExt.([]interface{})
	ModulesNopecl := []string{"memcached", "imagick", "msgpack", "xdebug", "redis"}

	DockerModules := []string{"iconv", "pdo_mysql", "pdo_sqlite", "mysqli", "gd", "exif", "intl", "xsl",
		"json", "soap", "dom", "zip", "opcache", "xml", "mbstring",
		"bz2", "calendar", "ctype", "bcmatch", "mcrypt",
	}
	var Dockerfile string
	if PhpVersion[confYaml.From].distrib == "alpine" {
		Dockerfile = Alpine(phpModules, ModulesNopecl, DockerModules, PhpVersion, maintainer, confYaml)
	} else if PhpVersion[confYaml.From].distrib == "debian" {
		Dockerfile = Debian(phpModules, ModulesNopecl, DockerModules, PhpVersion, maintainer, confYaml)
	}

	var myphp bool
	myphp = false
	if myphp {
		Dockerfile += "ADD php.ini /etc/php.ini" // use users php file
	}

	err = GenerateRunScript()
	if err != nil {
		return
	}

	ioutil.WriteFile("Dockerfile", []byte(Dockerfile), 0644)
}
