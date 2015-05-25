* split worker into libworker and cmdworker: the server can then be started with an embedded worker for simple installations, or delegate to a queue and a fleet of worker
* do we really need multiple executors per worker ? coudn't we start one worker per required executor ?
* global scm key neds to be reworked to be stored in the DB and expose an api endpoint to retrieve it
