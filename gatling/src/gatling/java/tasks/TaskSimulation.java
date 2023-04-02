package tasks;

import io.gatling.javaapi.core.ScenarioBuilder;
import io.gatling.javaapi.core.Simulation;
import io.gatling.javaapi.http.HttpProtocolBuilder;

import static io.gatling.javaapi.core.CoreDsl.rampUsers;
import static io.gatling.javaapi.core.CoreDsl.scenario;
import static io.gatling.javaapi.http.HttpDsl.http;

public class TaskSimulation extends Simulation {

    UserRequests user = new UserRequests();
    TaskRequests task = new TaskRequests();

    ScenarioBuilder active = scenario("new users").exec(
            user.create,
            user.login,
            task.create,
            task.update,
            task.delete,
            user.logout
    );

    HttpProtocolBuilder httpProtocol = http.baseUrl("http://localhost:1323")
            .contentTypeHeader("application/json");

    {
        setUp(
                active.injectOpen(rampUsers(30).during(30))
        ).protocols(httpProtocol);
    }
}
