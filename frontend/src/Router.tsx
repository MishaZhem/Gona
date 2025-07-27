import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import { Login, Main, Profile } from './pages'
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
                    {
                        path: '/profile',
                        element: (
                            <Profile />
                        ),
                    },
                ],
            },
        ],
    )

    return <RouterProvider router={routing} />
}
