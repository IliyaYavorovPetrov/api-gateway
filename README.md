# Problem definition 
At the beginning of many small projects adopt a microservices architecture, it is common practice to expose each service directly to the external world. This approach is straightforward because there is typically minimal load and complexity. However, as the microservices architecture scales, it becomes advisable to avoid direct exposure of most services to the external world. Doing so can help mitigate security vulnerabilities, simplify tracking and monitoring, and enable more effective implementation of rate limiting and performance management.
# Solution overview
The best solution in the case is to create a new service (API Gateway), which job is to handle all the incoming traffic to the system and take care of:
- authentication and authorization
- routing to the proper service/s
- logging and monitoring
- rate limiting
- transform requests into a more appropriate format to enhance performance

More information [[#Design|here]]
# In depth solution
## Tech stack
**Golang:**
- Golang is the primary language chosen for its exceptional performance, extensive documentation, and strong community support.
- It seamlessly integrates with Kubernetes, which is commonly used for microservices orchestration.

**Redis:**
- Very high performance and popular choice for session store
- Good client library for Go and big community 

**gRPC:**
- It enables quick and reliable data sharing among microservices, enhancing overall performance.
- Good integration with Go
- A user-friendly RESTful interface is retained for external interactions, while gRPC is used for high-speed backend operations, ensuring top-notch performance.

## Design

**Design Choices:**
- **Dynamic Routing:** We opted not to rely on a predefined configuration file for routing. Instead, routing information is stored directly in an Redis cluster. This approach allows us to update routing information in real-time without requiring a new API gateway deployment, enhancing system availability.

**Main Features:**
- **Authentication and Authorization:** Ensuring secure access by authenticating and authorizing incoming requests using sessions and cookies.
- **Dynamic Routing:** Utilizing Redis-stored routing information to direct requests to the appropriate service(s) without the need for frequent gateway deployments.
- **Rate Limiting:** Implementing rate limiting based on the IP address, with rate limit data also stored in the Redis cluster.
- **Logging and Statistics:** Gathering comprehensive statistics and logging for each request, including the count of different types of requests, network load volume, the number of rate-limited users, and latency metrics.
- **Request Transformation:** Transforming incoming RESTful requests into gRPC, facilitating communication with backend services.

These design choices and features collectively enhance the scalability, maintainability, and flexibility of our API gateway, ensuring robust and efficient management of incoming traffic.

# Problems with the solution
We must maintain awareness of the technical implementation of this service, as it can become a single point of failure for the entire system. Therefore, it should undergo rigorous testing, including unit, integration, and load tests.