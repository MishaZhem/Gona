import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import { Login, Main, Profile } from './pages'
import Template from './pages/Template/Template'
import ProtectedRoute from './AuthorizedRoute'

export default function Router() {
    const routing = createBrowserRouter(
        [
            {
                path: '/',
                element: <Template />,
                children: [
                    {
                        index: true,
                        element: (
                            <Main />
                        ),
                    },
                    {
                        path: 'login',
                        element: (
                            <Login />
                        ),
                    },
                    {
                        element: <ProtectedRoute />,
                        children: [
                            {
                                path: 'profile',
                                element: (
                                    <Profile />
                                ),
                            },
                        ],
                    },
                ],
            },
        ],
    )

    return <RouterProvider router={routing} />
}
