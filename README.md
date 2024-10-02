# gh-pull-request-dashboard
This is a dashboard to see the open Pull Requests for a given repository 

## configuration
Information on how to set everything up. In their current form these steps expect you 

first you should run ``npm i`` in the root of the project

### backend
all backend files live in ``backend/``.

The most important file here is in ``./backend/db/`` named ``example_config.json``. Copy this file and rename it to be just ``config.json`` Then fill in the fields with the following information:

```
{
  "token": "", -- Your GitHub Personal Access Token
  "owner": "google", -- Owner of the repository
  "repo": "go-github" -- Name of the repository
}
```

in ``backend/`` run ``go get .`` and ``go run .`` to run the backend

### frontend
navigate to ``frontend/`` and run ``npm i``

If you want to change the port the web-server runs on you can change the ``port`` values in ``vite.config.js``

If you are hosting the site on something other than a root domain, say ``example.com/pr-dashboard/`` for example then you must change the ``base`` field in ``vite.config.js`` and in ``svelte.config.js`` to be the same as the base you are running it in.

you can run or build the frontend from the root by running ``npm run dev`` or ``npm run build``.


## running the application

From the root directory you can either run ``npm run start`` which will run a development version of the application.


to run the built versions of the application it is required to build the backend with ``cd backend/`` & ``go build .`` and the frontend to have been built with ``npm run build``.

on Linux you can run with ``npm run lin`` to run a built version of the application. 

on Windows you can run ``npm run win`` to run a built version of the application, 


