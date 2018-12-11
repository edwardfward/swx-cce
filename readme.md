## SOFWERX CCE

### Stack Description 

* Front End - HTML5, JavaScript
* Back End - Go, NGINX, PostgreSQL, Docker, Kubernetes, Linux/OS

### User Category Definitions

* **Site Admin** - Administers the site, adds and removes admins, configuring secrets, and monitors site status and health.

* **Admin** - Responsible for adding removing CCE events and managing event participants.

* **Facilitator** - Runs the CCE, manages and presents results, edits data, sets and remove timers, and determines group participants.

* **Analyst** - Full access to CCE data, minus user data protected by privacy guidelines, to export and analyze.

* **User** - Participates in CCEs they have been invited to by submitting limitations, feedback, surveys, and other data to the facilitator. 


### Draft Use Cases

1. User opens a web browser on a desktop or mobile device and goes to the `cce.sofwerx.org` app to create a user account or authenticate with OAuth2, username and password, or via a one-time magic link sent to the user's validated email account.

2. Admin opens a web browser and goes to `cce.sofwerx.org/admin` page and logs in via OAuth2, username and password, or via a one-time magic link sent to the admin's @sofwerx.org account. 

3. Admin goes to `cce.sofwerx.org/admin` to create a new CCE event. Event needs to have a short, easy-to-remember title, start and end dates, location, and the ability to add and remove users from the CCE event before, during, and after the event takes place. The event title will be used throughout.

4. User receives a link to CCE via email or invite and clicks a link to `cce.sofwerx.org/{cce title}`. If the user's secure JSON Web Token ([JWT](https://jwt.io/)) is invalid, the user will be directed to the logon page to create an account or authenticate. Once the user successfully authenticates, they will be redirected to the original `cce.sofwerx.org/{cce title}` page, but only if the admin has invited or added the user to the CCE.

5. Facilitator goes to `cce.sofwerx.org/run/{cce title}` to view the CCE's facilitator dashboard for the event. Facilitators must be designated by an admin and authenticate via valid JWT, OAuth2, username and password, or one-time magic link sent to the facilitator's email.

### References

* [The Twelve-Factor App](https://12factor.net/)
* [JSON Schemas](http://json-schema.org/)
* [JSON Web Tokens](https://jwt.io/)
* [Go JSON Schema](https://github.com/xeipuuv/gojsonschema)
* [NGINX Ingress Controller for Kubernetes](https://www.nginx.com/products/nginx/kubernetes-ingress-controller)
* [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)  
* [Go Kit](https://gokit.io/)
* [Go Gorilla Websocket](https://github.com/gorilla/websocket)
* [OpenAPI](https://github.com/OAI/OpenAPI-Specification)

