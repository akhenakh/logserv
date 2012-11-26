logserv
=======

Expose log files via HTTP in a minute safely



config
======

{
   "port":8080,
   "auth_file":"htpasswd",
   "logfiles":[
      {"path":"/var/log/system.log", "users":[ "bob", "kevin"]},
	  {"path":"/var/log/install.log"}]
}



TODO
====

Detect behind proxy usage ? no websocket