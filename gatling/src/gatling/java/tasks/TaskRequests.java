package tasks;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.json.JsonMapper;
import io.gatling.javaapi.core.ChainBuilder;
import tasks.model.Task;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.http;
import static io.gatling.javaapi.http.HttpDsl.status;
import static tasks.model.Task.generateTask;

public class TaskRequests {

    JsonMapper mapper = new JsonMapper();

    public ChainBuilder create =
            exec(session -> {
                Task task = generateTask();
                String serializedTask = null;
                System.out.println("generated task " + "id" + task.getId() + "description" + task.getDescription() + "completed" + task.getCompleted() );
                try {
                    serializedTask = mapper.writeValueAsString(task);
                } catch (JsonProcessingException e) {
                    throw new RuntimeException(e);
                }
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
                Integer createdTaskId = session.getInt("createdTaskId");
                Task task = generateTask();
                String serializedTask = null;
                try {
                    serializedTask = mapper.writeValueAsString(task);
                } catch (JsonProcessingException e) {
                    throw new RuntimeException(e);
                }
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
