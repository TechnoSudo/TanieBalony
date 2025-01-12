package com.taniebalony.backendtest;

import org.springframework.http.HttpEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class PingController {
    @GetMapping
    String ping() {
        return "ping";
    }
    @GetMapping("/api")
    String ping2(HttpEntity<String> httpEntity) {
        String json = httpEntity.getBody();
        return "ping2";
    }
    @GetMapping("/api/ping")
    String ping3() {
        return "ping3";
    }
}
