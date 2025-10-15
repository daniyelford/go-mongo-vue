<?php

namespace App\Http;

use Illuminate\Foundation\Http\Kernel as HttpKernel;

class Kernel extends HttpKernel
{
    /**
     * Global HTTP middleware stack.
     */
    protected $middleware = [
        // middlewareهای سراسری
    ];

    /**
     * Route middleware groups.
     */
    protected $middlewareGroups = [
        'web' => [
            // web middlewareها
        ],

        'api' => [
            // api middlewareها
            'throttle:api',
            \Illuminate\Routing\Middleware\SubstituteBindings::class,
        ],
    ];

    /**
     * Route middleware.
     */
    protected $routeMiddleware = [
        'jwt.auth' => \App\Http\Middleware\JWTAuthMiddleware::class,
        // سایر middlewareها
    ];
}
