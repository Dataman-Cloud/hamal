### Install

#### Opt 1: in docker way

>```bash
docker run -it --net=host --rm -v /root/wtzhou/hamal/frontend:/src  mkenney/npm:latest sh
npm install
bower install --allow-root
gulp serve
```

#### Opt 2: in bare metal way

##### Install required tools `gulp` and `bower`:
```
npm install -g gulp bower
```

##### Go into frontend directory, and run:
```
npm install
```

```
bower install
```

## use Gulp tasks

* `gulp` or `gulp build` to build an optimized version of your application in `/dist`
* `gulp serve` to launch a browser sync server on your source files
* `gulp serve:dist` to launch a server on your optimized application
* `gulp test` to launch your unit tests with Karma
* `gulp test:auto` to launch your unit tests with Karma in watch mode
* `gulp protractor` to launch your e2e tests with Protractor
* `gulp protractor:dist` to launch your e2e tests with Protractor on the dist files
