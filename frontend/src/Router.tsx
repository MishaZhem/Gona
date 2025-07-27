import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import { Login, Main } from './pages'
import Template from './Template'

export default function Router() {
    const routing = createBrowserRouter(
        [
            {
                path: '/',
                element: <Template />,
                children: [
                    {
                        path: '/',
                        element: (
                            <Main />
                        ),
                    },
                    {
                        path: '/login',
                        element: (
                            <Login />
                        ),
                    },
                ],
            },
        ],
    )

    return <RouterProvider router={routing} />
}
