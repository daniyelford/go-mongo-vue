<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Firebase\JWT\JWT;
use Firebase\JWT\Key;
use Illuminate\Support\Facades\Redis;

class JWTAuthMiddleware
{
    public function handle(Request $request, Closure $next)
    {
        $authHeader = $request->header('Authorization');

        if (!$authHeader || !str_starts_with($authHeader, 'Bearer ')) {
            return response()->json(['message' => 'missing token'], 401);
        }

        $tokenString = substr($authHeader, 7); // remove 'Bearer '

        try {
            // decode JWT
            $credentials = JWT::decode($tokenString, new Key(env('JWT_SECRET'), 'HS256'));

            if (empty($credentials->mobile)) {
                return response()->json(['message' => 'invalid token claims'], 401);
            }

            $mobile = $credentials->mobile;

            // check token in Redis
            $val = Redis::get("token:$mobile");
            if (!$val || $val !== $tokenString) {
                return response()->json(['message' => 'expired or revoked token'], 401);
            }

            // store mobile in request attributes for controller
            $request->attributes->set('mobile', $mobile);

        } catch (\Exception $e) {
            return response()->json([
                'message' => 'invalid token',
                'error' => $e->getMessage()
            ], 401);
        }

        return $next($request);
    }
}
