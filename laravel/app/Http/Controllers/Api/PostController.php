<?php
namespace App\Http\Controllers\Api;
use App\Http\Controllers\GoApiController; // Base controller که sendToGo داره
use Illuminate\Http\Request;
class PostController extends GoApiController
{
    protected $routes = [
        'getAll' => '/api/posts/all',
        'create' => '/api/posts/create',
        'edit' => '/api/posts/edit',
        'delete' => '/api/posts/delete',
    ];
    public function call(Request $request, $action)
    {
        if (!isset($this->routes[$action])) {
            return response()->json(['status' => 'error', 'message' => 'متد نامعتبر']);
        }
        $endpoint = $this->routes[$action];

        return $this->sendToGo($endpoint, $request->all());
    }
}
