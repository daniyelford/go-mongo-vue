<?php
namespace App\Http\Controllers\Api;
use App\Http\Controllers\GoApiController; // Base controller که sendToGo داره
use Illuminate\Http\Request;
class PostController extends GoApiController
{
    protected $routes = [
        'getAll' => '/posts/all',
        'create' => '/posts/create',
        'edit' => '/posts/edit',
        'delete' => '/posts/delete',
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
