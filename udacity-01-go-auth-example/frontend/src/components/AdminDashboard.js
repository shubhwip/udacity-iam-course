import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

function AdminDashboard() {
  const [users, setUsers] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await axios.get('http://localhost:8080/admin/users', {
          headers: { Authorization: localStorage.getItem('token') }
        });
        setUsers(response.data);
      } catch (error) {
        console.error('Failed to fetch users:', error);
        navigate('/login');
      }
    };

    fetchUsers();
  }, [navigate]);

  return (
    <div>
      <h2>Admin Dashboard</h2>
      <ul>
        {users.map(user => (
          <li key={user.id}>{user.username} - {user.role}</li>
        ))}
      </ul>
    </div>
  );
}

export default AdminDashboard;
