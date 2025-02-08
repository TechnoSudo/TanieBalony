package com.taniebalony.gateway;

import java.security.Principal;

import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.factory.AbstractGatewayFilterFactory;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;

@Component
public class AuthFilter extends AbstractGatewayFilterFactory<AuthFilter.Config> {

    public AuthFilter(WebClient.Builder webClientBuilder) {
        super(Config.class);
    }

    @Override
    public GatewayFilter apply(Config config) {
        return (exchange, chain) -> exchange.getPrincipal()
                .map(Principal::getName)
                .defaultIfEmpty("Default User")
                .map(userName -> {
                    //adds header to proxied request
                    System.out.println("Config value ="+config.value);
                    exchange.getRequest().mutate().header(config.value,
                            userName).build();
                    System.out.println("Config First pre header filter" +
                            exchange.getRequest().getHeaders());
                    return exchange;
                })
                .flatMap(chain::filter);
    }

    @Override
    public Config newConfig() {
        return new Config();
    }

    public static class Config {
        private String value;

        public String getValue() {
            return value;
        }

        public void setValue(String value) {
            this.value = value;
        }
    }
}
