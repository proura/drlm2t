DRLM2T Certificate Generation
=============================

If you need to re-generate the certificate or want to extend expiration time
please re-run the following command with proper options:

::

  openssl req -newkey rsa:4096 -nodes -keyout ./drlm2t.key -x509 -days 1825 -subj "/C=ES/ST=CAT/L=GI/O=SA/CN=$(hostname -s)" -out ./drlm2t.crt