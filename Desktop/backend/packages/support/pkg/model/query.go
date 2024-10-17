package model

const (
	JoinSelect          string = "chat.sessions.id, chat.sessions.chat_type, chat.sessions.created, chat.session_statuses.status, chat.session_users.user_id, chat.session_read_receipts.last_read, chat.session_descriptors.name, chat.session_descriptors.value, chat.session_last_messages.last_message, chat.session_last_messages.last_sent, chat.session_patients.patient_id, chat.patients.name as patient_name, chat.patients.last_name as patient_last_name, chat.patients.address as patient_address, chat.patients.insurance_id as patient_insurance_id, chat.patients.birthday as patient_birthday"
	JoinSelectWithNotes string = "chat.sessions.id, chat.sessions.chat_type, chat.sessions.created, chat.session_statuses.status, chat.session_users.user_id, chat.session_read_receipts.last_read, chat.session_descriptors.name, chat.session_descriptors.value, chat.session_last_messages.last_message, chat.session_last_messages.last_sent, chat.session_notes.notes, chat.session_patients.patient_id, chat.patients.name as patient_name, chat.patients.last_name as patient_last_name, chat.patients.address as patient_address, chat.patients.insurance_id as patient_insurance_id, chat.patients.birthday as patient_birthday"
	JoinStatus          string = "INNER JOIN chat.session_statuses ON chat.session_statuses.session_id = chat.sessions.id"
	JoinUsers           string = "INNER JOIN chat.session_users ON chat.session_users.session_id = chat.sessions.id"
	JoinDescriptors     string = "LEFT JOIN chat.session_descriptors ON chat.session_descriptors.session_id = chat.sessions.id"
	JoinPatientsSession string = "LEFT JOIN chat.session_patients ON chat.session_patients.session_id = chat.sessions.id"
	JoinPatients        string = "LEFT JOIN chat.patients ON chat.patients.id = chat.session_patients.patient_id"
	JoinNotes           string = "LEFT JOIN chat.session_notes ON chat.session_notes.session_id = chat.sessions.id"
	JoinMessages        string = "LEFT JOIN chat.session_last_messages ON chat.session_last_messages.session_user_id = chat.session_users.id"
	JoinRead            string = "LEFT JOIN chat.session_read_receipts ON chat.session_read_receipts.session_user_id = chat.session_users.id"
	WhereSession        string = "chat.sessions.id = @id"
	OrderSent           string = "chat.session_last_messages.last_sent DESC"
	OrderCreatedSent    string = "chat.sessions.created, chat.session_last_messages.last_sent DESC"
	WhereIdByUser       string = "chat.sessions.id IN (SELECT chat.session_users.session_id FROM chat.session_users WHERE chat.session_users.user_id = @id)"
	LastMessageUpsert   string = "INSERT INTO chat.session_last_messages (session_user_id, last_message, last_sent) SELECT id, @message, @sent FROM chat.session_users WHERE session_id = @id AND user_id = @user ON CONFLICT (session_user_id) DO UPDATE SET last_message = EXCLUDED.last_message, last_sent = EXCLUDED.last_sent"
	LastReadUpsert      string = "INSERT INTO chat.session_read_receipts (session_user_id, last_read) SELECT id, @read FROM chat.session_users WHERE session_id = @id AND user_id = @user ON CONFLICT (session_user_id) DO UPDATE SET last_read = EXCLUDED.last_read"
	WhereRole           string = "chat.sessions.chat_type = @chat_type AND (chat.session_statuses.status = 'PENDING' OR " + WhereIdByUser + ")"
)
