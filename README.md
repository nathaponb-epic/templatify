# Templatify

The cli based application for manupulate the ADMD templates content.

## Download and Installation

### Windows
* Download the latest version .exe file [templatify.exe]("./bin/windows/templatify.exe")
* Place it inside a directory of choice
* Update *PATH Variable* to templatify directory on your System

### macOS
* Download the latest version binary file [templatify]("./bin/macos/templatify")

### Linux
* Download the latest version binary file [templatify]("./bin/linux/templatify")
---


## Commands
Templatify version 0.0.1
* **cdnify**: *for update all related paths in the template to extenal content delivery service (CDN)*.
* **localify**: *for update all the related paths back to local assets*

---
## Configuration
templatify.yaml is the **mandatory** config file for templatify CLI, 
place it in the target root directory of your template with the being executed.
structure in YAML syntax.

### structure

* Top-level key: `commands:`
* Command: Each command enclosed in a dictionary yaml data type which contains necessary attributes `name, domain, path, ...`
    * `name`: the name of command [see available commands on release note](#)
    * `domain`: the domain name (*in case of localify leave this field blank*)
    * `path`: url path
    * `app_prefix`: prefix name to be replace the **PREFIX** variable in *prefix.js* which being refered to mainly in .js file.
    * `root_path`: a dictionary of root path specifically to each catagory of file. [see more detail of root_path](#attr_root_path)
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
---

## <a id="attr_root_path"></a>root_path



<details>
  <summary>How to know what is the root path of each type of catagory of file ?</summary>
  
  Basically you have to grasp how each template is structured, e.g. ncp template
  ```
    |--README.md
    |--default
        |--fonts
        |--images
        |--mail
        |--script
            |--constant
            |--css
            |--js
        |--views
  ```
  note: *default directory is excluded.* 

  **root_path** attribute

  * `font`: /fonts
  * `css`: /script
  * `script`: /scripts
  * `constant`: /scripts
  * `image`: /images

  
</details>

