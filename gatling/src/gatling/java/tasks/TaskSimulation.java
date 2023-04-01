package tasks;

import io.gatling.javaapi.core.ScenarioBuilder;
import io.gatling.javaapi.core.Simulation;
import io.gatling.javaapi.http.HttpProtocolBuilder;

import static io.gatling.javaapi.core.CoreDsl.rampUsers;
import static io.gatling.javaapi.core.CoreDsl.scenario;
import static io.gatling.javaapi.http.HttpDsl.http;

public class TaskSimulation extends Simulation {

    UserRequests user = new UserRequests();

    ScenarioBuilder active = scenario("active users").exec(user.create, user.login, user.logout);

//    ScenarioBuilder lazy = scenario("active users").exec(search, browse);

    HttpProtocolBuilder httpProtocol = http.baseUrl("http://localhost:1323")
            .contentTypeHeader("application/json");

    {
        setUp(
                active.injectOpen(rampUsers(5).during(10))
        ).protocols(httpProtocol);
    }
}
