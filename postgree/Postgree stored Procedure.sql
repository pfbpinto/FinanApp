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
    IN p_ParentIncomeAmount NUMERIC(15,2),
    IN p_BeginDate DATE,
    OUT p_Message TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    NewFinancialUserItemid INT;
    v_CurrentDate DATE;
    v_Iterations INT;
    v_Increment INTERVAL;
    i INT;
BEGIN
    -- Validações de campos obrigatórios
    IF p_UserID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserID"}';
        RETURN;
    END IF;
    IF p_FinancialUserItemName IS NULL OR p_FinancialUserItemName = '' THEN
        p_Message := '{"status": "fail", "message": "Missing FinancialUserItemName"}';
        RETURN;
    END IF;
    IF p_RecurrencyID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing RecurrencyID"}';
        RETURN;
    END IF;
    IF p_FinancialUserEntityItemID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing FinancialUserEntityItemID"}';
        RETURN;
    END IF;
    IF p_ParentIncomeAmount IS NULL OR p_ParentIncomeAmount <= 0 THEN
        p_Message := '{"status": "fail", "message": "Invalid ParentIncomeAmount"}';
        RETURN;
    END IF;
    IF p_BeginDate IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing BeginDate"}';
        RETURN;
    END IF;

    -- Inicia transação manualmente
    BEGIN
        -- Insere um novo Parent Income
        INSERT INTO FinancialUserItem (
            FinancialUserItemName, EntityID, UserEntityID, RecurrencyID,
            FinancialUserEntityItemID, ParentFinancialUserItemID
        ) VALUES (
            p_FinancialUserItemName, 5, p_UserID, p_RecurrencyID,
            p_FinancialUserEntityItemID, NULL
        ) RETURNING FinancialUserItemID INTO NewFinancialUserItemid;

        -- Garante que a inserção foi bem-sucedida
        IF NewFinancialUserItemid IS NULL THEN
            RAISE EXCEPTION 'Error inserting FinancialUserItem';
        END IF;

        -- Configura número de inserções e intervalo conforme a recorrência
        IF p_RecurrencyID = 1 THEN -- One time
            v_Iterations := 1;
            v_Increment := '1 day'::INTERVAL;
        ELSIF p_RecurrencyID = 2 THEN -- Monthly
            v_Iterations := 12;
            v_Increment := '1 month'::INTERVAL;
        ELSIF p_RecurrencyID = 3 THEN -- Quarterly
            v_Iterations := 4;
            v_Increment := '4 months'::INTERVAL;
        ELSIF p_RecurrencyID = 4 THEN -- Yearly
            v_Iterations := 1;
            v_Increment := '1 year'::INTERVAL;
        ELSE
            RAISE EXCEPTION 'Invalid RecurrencyID';
        END IF;

        -- Insere múltiplos registros conforme a recorrência
        v_CurrentDate := p_BeginDate;
        FOR i IN 1..v_Iterations LOOP
            INSERT INTO UserFinancialForecast (
                usercategoryid, financialuseritemid, userfinancialforecastbegindate,
                userfinancialforecastenddate, userfinancialforecastamount, currencyid
            ) VALUES (
                NULL, NewFinancialUserItemid, v_CurrentDate,
                CASE WHEN p_RecurrencyID IN (1, 4) THEN v_CurrentDate + v_Increment - '1 day'::INTERVAL ELSE NULL END,
                p_ParentIncomeAmount, 1
            );
           
            -- Atualiza data para a próxima recorrência
            v_CurrentDate := v_CurrentDate + v_Increment;
        END LOOP;

        -- Se tudo deu certo, define mensagem de sucesso
        p_Message := '{"status": "success", "message": "New user income created successfully."}';

    EXCEPTION
        WHEN OTHERS THEN
            -- Captura erro e define mensagem de falha sem tentar rollback manualmente
            p_Message := '{"status": "fail", "message": "Error creating forecast values: ' || SQLERRM || '"}';
    END;
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