package tasks;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.json.JsonMapper;
import tasks.model.User;
import io.gatling.javaapi.core.ChainBuilder;
import org.apache.commons.lang3.RandomStringUtils;

import java.util.Collections;
import java.util.Iterator;
import java.util.Map;
import java.util.function.Supplier;
import java.util.stream.Stream;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.http;
import static io.gatling.javaapi.http.HttpDsl.status;
import static tasks.model.User.generateRandomUser;

public class UserRequests {
    JsonMapper mapper = new JsonMapper();
    Iterator<Map<String, Object>> userFeeder =
            Stream.generate((Supplier<Map<String, Object>>) () -> {
                        User u = generateRandomUser();
                        String serializedUser = null;
                        try {
                            serializedUser = mapper.writeValueAsString(u);
                        } catch (JsonProcessingException e) {
                            throw new RuntimeException(e);
                        }
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
                User generatedUser = null;
                try {
                    generatedUser = mapper.readValue(generatedUserBody, User.class);
                } catch (JsonProcessingException e) {
                    throw new RuntimeException(e);
                }
                User loginUser = new User(null, generatedUser.getEmail(), generatedUser.getPassword());
                String loginUserSerialized = null;
                try {
                    loginUserSerialized = mapper.writeValueAsString(loginUser);
                } catch (JsonProcessingException e) {
                    throw new RuntimeException(e);
                }
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
