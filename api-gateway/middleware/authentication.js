// Create the authentication middleware object
var authentication = new TykJS.TykMiddleware.NewMiddleware({});

// Middleware function to validate JWT from the Authorization header
authentication.NewProcessRequest(function(request, session, spec) {
    var authHeader = request.Headers["Authorization"];
    
    // Ensure authHeader is a string
    if (Array.isArray(authHeader)) {
        authHeader = authHeader[0];
    }
    if (!authHeader) {
        console.log("Missing Authorization header");
        return authentication.ReturnError("Missing Authorization header", 401);
    }
    authHeader = String(authHeader);
    
    // Expect header in "Bearer <token>" format
    var parts = authHeader.split(" ");
    if (parts.length !== 2 || parts[0] !== "Bearer") {
        console.log("Invalid Authorization header format");
        return authentication.ReturnError("Invalid Authorization header format", 401);
    }

    var token = parts[1];
    console.log("Extracted JWT:", token);

    // Split the JWT into parts
    var jwtParts = token.split(".");
    if (jwtParts.length !== 3) {
        console.log("Invalid JWT format");
        return authentication.ReturnError("Invalid JWT token", 401);
    }

    // Decode JWT payload (second part)
    var payloadJson = TykJS.Base64Decode(jwtParts[1]);
    var payload;
    try {
        payload = JSON.parse(payloadJson);
    } catch (e) {
        console.log("Error parsing JWT payload: " + e);
        return authentication.ReturnError("Invalid JWT payload", 401);
    }

    // Check if the JWT has expired
    var currentTimestamp = Math.floor(Date.now() / 1000);
    if (!payload.exp || payload.exp < currentTimestamp) {
        console.log("JWT is expired. Expiration:", payload.exp, "Current time:", currentTimestamp);
        return authentication.ReturnError("JWT expired", 401);
    }

    // If valid, pass the request along
    return authentication.ReturnData(request, session.meta_data || {});
});

// Register the middleware with Tyk
TykJS.TykMiddleware.AddMiddleware("authentication", authentication);
