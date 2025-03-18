function checkHealth(request, session, config) {
    var services = {
        "account": "http://account-service:3000",
        "backup": "http://backup-service:3001",
        "dashboard": "http://dashboard-service:3002"
    };

    var results = {};
    var allHealthy = true;

    for (var service in services) {
        var url = services[service];
        
        var response = TykMakeHttpRequest({
            Method: "GET",
            Headers: {
                "Content-Type": "application/json"
            },
            Domain: url
        });

        if (response.Code === 200) {
            try {
                var body = JSON.parse(response.Body);
                results[service] = body;
                if (body.status !== "success") {
                    allHealthy = false;
                }
            } catch (e) {
                results[service] = { status: "error", error: "Invalid JSON response" };
                allHealthy = false;
            }
        } else {
            results[service] = { status: "unreachable", error: response.Body };
            allHealthy = false;
        }
    }

    var finalResponse = {
        status: allHealthy ? "healthy" : "degraded",
        services: results
    };

    return TykJsResponse({
        Body: JSON.stringify(finalResponse),
        Headers: {
            "Content-Type": "application/json"
        },
        Code: 200
    });
}
