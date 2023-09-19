# Templatify

the cli based application for updating the ADMD templates

### templatify.yaml
templatify.yaml is the **mandatory** config file for templatify CLI to run properly, place it in the target root directory of your template with the being executed.

## Commands
Templatify version 0.0.1 Beta
* cdnify: *for update all related paths in the template to extenal content delivery service (CDN)*
* localify: *for update all the related paths back to local assets*

---

structure in YAML syntax

* Top-level key: `commands:`
* Command dictionary: `name, domain, path, ...`
    * `name`: the name of command
    * `domain`: the domain name (*in case of localify leave this field blank*)
    * `path`: url path
    * `app_prefix`: prefix name to be replace the **PREFIX** variable in *prefix.js*
    * `root_path`: a dictionary of root path specifically to each catagory of file
    * `ignore_dir`: a list of directories to be ignored processing (*typically public libs, images, fonts, contants, mail*)
    * `ignore_file`: a list of files to be ignored processing (*typically public libs file e.g. bootstrap.min.css, config files e.g. licence file and any file that has nothing to do with changing path*)

```yml
commands:
    - name: cdnify
      domain: https://cdn.delivery.com
      path: /ncp
      app_prefix: ncpCDN
      root_path:
        image: /images
        css: /script
        script: /script
        constant: /script
        font: /fonts
      ignore_dir:
        - fonts
        - images
        - mail
        - constant
        - font-awesome-4.7.0
        - font-awesome-6.2.1
        - sweetalert2-8.17.4
      ignore_file:
        - bootstrap-datepicker.min.css
        - bootstrap.min.css
        - bootstrap.min.js
        - bs-datepicker.js
        - popper.min.js
        - jsrsasign-all-min.js

    - name: localify
      domain:
      path: /auth
      app_prefix: ncp
      root_path:
        image: /images
        css: /script
        script: /script
        constant: /script
        font: /fonts
      ignore_dir:
        - fonts
        - images
        - mail
        - constant
        - font-awesome-4.7.0
        - font-awesome-6.2.1
        - sweetalert2-8.17.4
      ignore_file:
        - bootstrap-datepicker.min.css
        - bootstrap.min.css
        - bootstrap.min.js
        - bs-datepicker.js
        - popper.min.js
        - jsrsasign-all-min.js
    
```

