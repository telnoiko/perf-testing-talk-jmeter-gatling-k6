package tasks;

import io.gatling.javaapi.core.ChainBuilder;
import tasks.model.User;

import java.util.Collections;
import java.util.Iterator;
import java.util.Map;
import java.util.function.Supplier;
import java.util.stream.Stream;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.http;
import static io.gatling.javaapi.http.HttpDsl.status;
import static tasks.model.User.generateRandomUser;
import static tasks.model.Util.deserializeFromString;
import static tasks.model.Util.serializeToString;

public class UserRequests {
    Iterator<Map<String, Object>> userFeeder =
            Stream.generate((Supplier<Map<String, Object>>) () -> {
                        User u = generateRandomUser();
                        String serializedUser = serializeToString(u);
                        return Collections.singletonMap("user", serializedUser);
                    }
            ).iterator();

    public ChainBuilder create = feed(userFeeder)
            .exec(http("create user").post("/users")
                    .body(StringBody("#{user}"))
                    .check(status().is(201))
                    .check(bodyString().exists()
                            .saveAs("createdUserBody"))
            );

    public ChainBuilder login =
            exec(session -> {
                String generatedUserBody = session.getString("user");
                User generatedUser = deserializeFromString(generatedUserBody, User.class);
                User loginUser = new User(null, generatedUser.getEmail(), generatedUser.getPassword());
                String loginUserSerialized = serializeToString(loginUser);
                return session.set("loginUser", loginUserSerialized);
            }).exec(http("login").post("/users/login")
                    .body(StringBody("#{loginUser}"))
                    .check(status().is(200))
                    .check(jsonPath("$.token").exists()
                            .saveAs("token"))
            );

    public ChainBuilder logout =
            exec(http("logout").post("/users/logoutAll")
                    .header("authorization", "bearer #{token}")
                    .check(status().is(200))
            );
}
