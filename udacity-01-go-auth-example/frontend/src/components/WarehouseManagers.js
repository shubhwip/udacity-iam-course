import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

function WarehouseManagersDashboard() {
  const [warehouseManagers, setWarehouseManagers] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchWarehouseManagers = async () => {
      try {
        const response = await axios.get('http://localhost:8080/admin/warehousemanagers', {
          headers: { Authorization: localStorage.getItem('token') }
        });
        setWarehouseManagers(response.data);
      } catch (error) {
        console.error('Failed to fetch warehouse managers:', error);
        navigate('/login');
      }
    };

    fetchWarehouseManagers();
  }, [navigate]);

  return (
    <div>
      <h2>Warehouse Managers Dashboard</h2>
      <ul>
        {warehouseManagers.map(warehouseManager => (
          <li key={warehouseManager.id}>{warehouseManager.name} - {warehouseManager.contact}</li>
        ))}
      </ul>
    </div>
  );
}

export default WarehouseManagersDashboard;
