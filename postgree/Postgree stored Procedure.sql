/*
ALL ACTIVITIES THAT REQUIRES CREATION, DELETION OR UPDATE WILL BE HANDLED VIA STORED PROCEDURES.
THE STRUCTURE TO RECEIVE A RESPONSE FOR ALL PROCEDURES WILL BE THE SAME, A MESSAGE WILL CAPTURE THE DATA BASE VALIDATION.
*/ 
 
 --STORED PROCEDURE TO CREATE USERS 
 
CREATE OR REPLACE PROCEDURE CreateUser(
    IN p_FirstName VARCHAR(100),
    IN p_LastName VARCHAR(255),
	IN p_EmailAddress VARCHAR(255),
	IN p_UserPassword VARCHAR(150),
    IN p_DateOfBirth DATE,
    OUT p_Message TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    existing_user INT;
BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF TRIM(p_FirstName) = '' OR TRIM(p_LastName) = '' OR p_DateOfBirth IS NULL OR 
       TRIM(p_UserPassword) = '' OR TRIM(p_EmailAddress) = '' IS NULL THEN
        p_Message := '{"status": "fail", "message": "Todos os campos são obrigatórios."}';
        RETURN;
    END IF;

    -- Verifica se o e-mail já existe
    SELECT COUNT(*) INTO existing_user FROM UserProfile WHERE EmailAddress = p_EmailAddress;
    IF existing_user > 0 THEN
        p_Message := '{"status": "fail", "message": "O e-mail já está cadastrado."}';
        RETURN;
    END IF;

    -- Insere o novo usuário
    INSERT INTO UserProfile (FirstName, LastName, DateOfBirth, UserPassword, EmailAddress, UserSubscription)
    VALUES (p_FirstName, p_LastName, p_DateOfBirth, p_UserPassword, p_EmailAddress, '0');

    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "Usuário criado com sucesso."}';
END;
$$;

-- STORED PROCEDURE TO CREATE NEW PARENT INCOME

CREATE OR REPLACE PROCEDURE CreateUserParentIncome(
    IN p_UserID INT,
    IN p_FinancialUserItemName VARCHAR(255),
	IN p_RecurrencyID INT,
	IN p_FinancialUserEntityItemID INT,
    OUT p_Message TEXT
)
LANGUAGE plpgsql
AS $$

BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF p_UserID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserProfileID}';
        RETURN;
    END IF;
	
    -- Insere o novo Parent Income
    INSERT INTO FinancialUserItem (
        FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, FinancialUserEntityItemID, ParentFinancialUserItemID
    ) VALUES (
        p_FinancialUserItemName, 5, p_UserID, p_RecurrencyID, p_FinancialUserEntityItemID, NULL
    );
    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "New user income created successfully."}';
END;
$$;



-- STORED PROCEDURE TO CREATE NEW PARENT INCOME TAX

CREATE OR REPLACE PROCEDURE CreateUserChildIncomeTax(
    IN p_UserID INT,
    IN p_FinancialUserItemName VARCHAR(255),
	IN p_RecurrencyID INT,
	IN p_FinancialUserEntityItemID INT,
	IN p_ParentFinancialUserItemID INT,
    OUT p_Message TEXT
)
LANGUAGE plpgsql
AS $$

BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF p_UserID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserProfileID}';
        RETURN;
    END IF;

    IF p_ParentFinancialUserItemID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Parent Income Item is mandatory}';
        RETURN;
    END IF;
    -- Insere o novo Child Tax Income
    INSERT INTO FinancialUserItem (
        FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, FinancialUserEntityItemID, ParentFinancialUserItemID
    ) VALUES (
        p_FinancialUserItemName, 7, p_UserID, p_RecurrencyID, p_FinancialUserEntityItemID, p_ParentFinancialUserItemID
    );
    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "New user tax associated to the parent income created successfully."}';
END;
$$;



-- STORED PROCEDURE TO CREATE NEW PARENT INCOME EXPENSE


CREATE OR REPLACE PROCEDURE CreateUserChildIncomeExpense(
    IN p_UserID INT,
    IN p_FinancialUserItemName VARCHAR(255),
	IN p_RecurrencyID INT,
	IN p_FinancialUserEntityItemID INT,
	IN p_ParentFinancialUserItemID INT,
    OUT p_Message TEXT
)
LANGUAGE plpgsql
AS $$

BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF p_UserID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserProfileID}';
        RETURN;
    END IF;

    IF p_ParentFinancialUserItemID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Parent Income Item is mandatory}';
        RETURN;
    END IF;
    -- Insere o novo Child Expense Income
    INSERT INTO FinancialUserItem (
        FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, FinancialUserEntityItemID, ParentFinancialUserItemID
    ) VALUES (
        p_FinancialUserItemName, 8, p_UserID, p_RecurrencyID, p_FinancialUserEntityItemID, p_ParentFinancialUserItemID
    );

    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "New user expense associated to the parent income created successfully."}';
END;
$$;