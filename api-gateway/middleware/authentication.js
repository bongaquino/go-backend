// Create the authentication middleware object
var authentication = new TykJS.TykMiddleware.NewMiddleware({});

// Middleware function to log JWT from Authorization header and pass the request
authentication.NewProcessRequest(function(request, session, spec) {
    console.log("Authentication middleware executed");

    // Extract the Authorization header from the request
    var authHeader = request.Headers["Authorization"];

    // Ensure authHeader is a string
    if (Array.isArray(authHeader)) {
        authHeader = authHeader[0];
    }
    if (authHeader === undefined || authHeader === null) {
        authHeader = "";
    } else {
        authHeader = String(authHeader);
    }

    // Extract the JWT by checking if "Bearer " exists
    var parts = authHeader.split(" ");
    var token = (parts.length === 2 && parts[0] === "Bearer") ? parts[1] : null;
    
    console.log("Extracted JWT:", token);

    // Pass the request
    return authentication.ReturnData(request, session.meta_data || {});
});

// Register the middleware with Tyk
TykJS.TykMiddleware.AddMiddleware("authentication", authentication);
