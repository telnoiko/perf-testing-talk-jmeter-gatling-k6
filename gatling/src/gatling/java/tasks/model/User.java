package tasks.model;

import org.apache.commons.lang3.RandomStringUtils;

public class User {
    public String name;
    public String email;
    public String password;

    public User() {
    }

    public User(String name, String email, String password) {
        this.name = name;
        this.email = email;
        this.password = password;
    }

    public static User generateRandomUser() {
        String email = RandomStringUtils.randomAlphabetic(3) + "@example.com";
        String password = RandomStringUtils.randomAlphabetic(8);
        String name = "name-" + RandomStringUtils.randomAlphabetic(3);
        return new User(name, email, password);
    }

    public String getName() {
        return name;
    }

    public String getEmail() {
        return email;
    }

    public String getPassword() {
        return password;
    }
}
