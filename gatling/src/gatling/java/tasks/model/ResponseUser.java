package tasks.model;

public class ResponseUser {
    User user;
    String token;

    public ResponseUser() {
    }

    public ResponseUser(User user, String token) {
        this.user = user;
        this.token = token;
    }
}
