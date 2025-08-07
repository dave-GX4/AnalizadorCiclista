import { useState, type FormEvent } from 'react';
import './App.css';

interface FormData {
  nombre: string;
  apellido_paterno: string;
  apellido_materno: string;
  email: string;
  sexo: 'M' | 'F' | '';
  categoria: string;
  pago_realizado: boolean;
  ine_path: string;
  comprobante_pago_path: string;
}

function App() {
  const [formData, setFormData] = useState<FormData>({
    nombre: '',
    apellido_paterno: '',
    apellido_materno: '',
    email: '',
    sexo: '',
    categoria: 'Aficionado',
    pago_realizado: true,
    ine_path: '',
    comprobante_pago_path: '',
  });

  const [isLoading, setIsLoading] = useState(false);
  const [serverMessage, setServerMessage] = useState<string | null>(null);
  const [isError, setIsError] = useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const isCheckbox = type === 'checkbox';
    if (isCheckbox) {
      const { checked } = e.target as HTMLInputElement;
      setFormData({ ...formData, [name]: checked });
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const transformToDSL = (data: FormData): string => {
    let dslString = '';
    dslString += `nombre: "${data.nombre}";\n`;
    dslString += `apellido_paterno: "${data.apellido_paterno}";\n`;
    if (data.apellido_materno) dslString += `apellido_materno: "${data.apellido_materno}";\n`;
    dslString += `email: "${data.email}";\n`;
    dslString += `sexo: "${data.sexo}";\n`;
    dslString += `categoria: "${data.categoria}";\n`;
    dslString += `pago_realizado: ${data.pago_realizado};\n`;
    if (data.ine_path) dslString += `ine_path: "${data.ine_path}";\n`;
    if (data.comprobante_pago_path) dslString += `comprobante_pago_path: "${data.comprobante_pago_path}";\n`;
    return dslString;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setServerMessage(null);
    setIsError(false);

    const dslPayload = transformToDSL(formData);
    
    try {
      const response = await fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: { 'Content-Type': 'text/plain' },
        body: dslPayload,
      });

      const result = await response.json();
      
      console.log("Respuesta del servidor:", result);

      if (!response.ok) {
        throw new Error(JSON.stringify(result.error) || 'Ocurrió un error en el servidor');
      }
      
      setIsError(false);
      
      // ==========================================================
      // === SECCIÓN MODIFICADA PARA UN MENSAJE MÁS AMIGABLE ===
      // ==========================================================
      if (result && result.participant) {
        // 1. Extraemos los datos del participante para que el código sea más limpio.
        const participant = result.participant;

        // 2. Construimos el nombre completo, manejando el apellido materno opcional.
        const nameParts = [participant.nombre, participant.apellido_paterno];
        if (participant.apellido_materno) {
          nameParts.push(participant.apellido_materno);
        }
        const fullName = nameParts.join(' ');
        
        // 3. Creamos el mensaje de éxito final.
        const successMessage = `¡Registro Exitoso!\n\nParticipante: ${fullName}\nCódigo de Registro: ${participant.participant_code}`;

        // 4. Establecemos el mensaje de éxito.
        setServerMessage(successMessage);
      } else {
        // Mensaje de respaldo si la respuesta no tiene la estructura esperada.
        setServerMessage("¡Registro exitoso! Revisa tu correo para ver los detalles.");
      }
      
    } catch (error: any) {
      setIsError(true);
      setServerMessage(`Error: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container">
      <h1>Registro de Ciclistas</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-grid">
          <input name="nombre" value={formData.nombre} onChange={handleInputChange} placeholder="Nombre(s)" required />
          <input name="apellido_paterno" value={formData.apellido_paterno} onChange={handleInputChange} placeholder="Apellido Paterno" required />
          <input name="apellido_materno" value={formData.apellido_materno} onChange={handleInputChange} placeholder="Apellido Materno" />
          <input name="email" type="email" value={formData.email} onChange={handleInputChange} placeholder="Correo Electrónico (Gmail)" required />
          <select name="sexo" value={formData.sexo} onChange={handleInputChange} required>
            <option value="" disabled>Selecciona tu sexo...</option>
            <option value="M">Masculino</option>
            <option value="F">Femenino</option>
          </select>
          <select name="categoria" value={formData.categoria} onChange={handleInputChange}>
            <option value="Elite">Elite</option>
            <option value="Aficionado">Aficionado</option>
            <option value="Juvenil">Juvenil</option>
          </select>
          <input name="ine_path" value={formData.ine_path} onChange={handleInputChange} placeholder="Ruta simulada INE (ej: /docs/ine.pdf)" />
          <input name="comprobante_pago_path" value={formData.comprobante_pago_path} onChange={handleInputChange} placeholder="Ruta simulada Pago (ej: /docs/pago.pdf)" />
        </div>
        <label className="checkbox-label">
          <input type="checkbox" name="pago_realizado" checked={formData.pago_realizado} onChange={handleInputChange} />
          Pago ya realizado
        </label>
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Registrando...' : 'Registrar Participante'}
        </button>
      </form>
      {serverMessage && (
        <div className={`server-message ${isError ? 'error' : 'success'}`}>
          <pre>{serverMessage}</pre>
        </div>
      )}
    </div>
  );
}

export default App;