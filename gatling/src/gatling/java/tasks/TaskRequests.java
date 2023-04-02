package tasks;

import io.gatling.javaapi.core.ChainBuilder;
import tasks.model.Task;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.http;
import static io.gatling.javaapi.http.HttpDsl.status;
import static tasks.model.Task.generateTask;
import static tasks.model.Util.serializeToString;

public class TaskRequests {

    public ChainBuilder create =
            exec(session -> {
                Task task = generateTask();
                String serializedTask = serializeToString(task);
                return session.set("newTask", serializedTask);
            })
                    .exec(http("create task").post("/tasks")
                            .body(StringBody("#{newTask}"))
                            .header("authorization", "bearer #{token}")
                            .check(status().is(201))
                            .check(bodyString().exists())
                            .check(jsonPath("$.id").exists()
                                    .saveAs("createdTaskId"))
                    );

    public ChainBuilder update =
            exec(session -> {
                Task task = generateTask();
                String serializedTask = serializeToString(task);
                return session.set("updatedTask", serializedTask);
            })
                    .exec(http("update task").put("/tasks/#{createdTaskId}")
                            .header("authorization", "bearer #{token}")
                            .body(StringBody("#{updatedTask}"))
                            .check(status().is(200))
                            .check(bodyString().exists()
                                    .saveAs("createdTaskBody"))
                    );

    public ChainBuilder delete =
            exec(http("delete task").delete("/tasks/#{createdTaskId}")
                    .header("authorization", "bearer #{token}")
                    .check(status().is(204))
            );
}
