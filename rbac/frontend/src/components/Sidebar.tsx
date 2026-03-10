import { NavLink } from 'react-router-dom';
import { usePermission } from '../auth/usePermission';

interface MenuItem {
  path: string;
  label: string;
  permission?: string;
}

const mainMenu: MenuItem[] = [
  { path: '/dashboard', label: 'Dashboard' },
  { path: '/products', label: 'Products', permission: 'products:read' },
  { path: '/orders', label: 'Orders', permission: 'orders:read' },
];

const adminMenu: MenuItem[] = [
  { path: '/users', label: 'Users', permission: 'users:read' },
  { path: '/roles', label: 'Roles', permission: 'roles:read' },
  { path: '/permissions', label: 'Permissions', permission: 'roles:read' },
];

export default function Sidebar() {
  const { hasPermission } = usePermission();

  const renderMenu = (items: MenuItem[]) =>
    items
      .filter((item) => !item.permission || hasPermission(item.permission))
      .map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={({ isActive }) =>
            `block px-4 py-2.5 rounded-lg text-sm font-medium transition-colors ${
              isActive
                ? 'bg-blue-100 text-blue-700'
                : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
            }`
          }
        >
          {item.label}
        </NavLink>
      ));

  const filteredAdmin = adminMenu.filter(
    (item) => !item.permission || hasPermission(item.permission)
  );

  return (
    <aside className="w-56 bg-white border-r border-gray-200 min-h-[calc(100vh-64px)] p-4 flex flex-col gap-1">
      {renderMenu(mainMenu)}
      {filteredAdmin.length > 0 && (
        <>
          <hr className="my-2 border-gray-200" />
          {renderMenu(adminMenu)}
        </>
      )}
    </aside>
  );
}
