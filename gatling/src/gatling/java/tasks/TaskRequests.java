package tasks;

import com.fasterxml.jackson.databind.json.JsonMapper;
import tasks.model.Task;
import io.gatling.javaapi.core.ChainBuilder;
import org.apache.commons.lang3.RandomStringUtils;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.http;
import static io.gatling.javaapi.http.HttpDsl.status;

public class TaskRequests {

    public ChainBuilder create = exec(session -> session.set("newTask", generateTask()))
            .exec(http("create task").post("/tasks")
                    .body(StringBody("#{newTask}"))
                    .check(status().is(201))
                    .check(bodyString().exists()
                            .saveAs("createdTaskBody"))
            );
    JsonMapper mapper = new JsonMapper();

    public Task generateTask() {
        String description = RandomStringUtils.randomAlphabetic(10);
        Boolean completed = Math.random() < 0.5;
        Task t = new Task(description, completed);
        return t;
    }
}
