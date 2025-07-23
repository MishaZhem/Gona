import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import { Main } from './pages'
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
                ],
            },
        ],
    )

    return <RouterProvider router={routing} />
}
